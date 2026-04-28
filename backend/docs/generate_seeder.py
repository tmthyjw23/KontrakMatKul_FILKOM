import random, re, os

BASE = os.path.dirname(os.path.abspath(__file__))
IN   = os.path.join(BASE, "extracted_all.txt")
OUT  = os.path.join(BASE, "migrations", "003_seed_data.sql")

# ── Lecturers ────────────────────────────────────────────────
LECTURERS = [
    "Andrew T. Liem","Green Mandias","Stenly R. Pungus","Debby E. Sondakh",
    "Edson Y. Putra","Green A. Sandag","Jacquline M. S. Waworundeng",
    "Jimmy H. Moedjahedy","Joe Y. Mambu","Lidya C. Laoh","Marchell Tombeng",
    "Oktoverano H. Lengkong","Reymon Rotikan","Reynoldus A. Sahulata",
    "Rolly Lontaan","Semmy W. Taju","Stenly I. Adam","Raissa Camila",
    "Andrew Tambunan","Andria K. Wahyudi","Wilsen Mokodaser",
]
N_LEC = len(LECTURERS)

# ── Parse extracted_all.txt ──────────────────────────────────
with open(IN, encoding="utf-8") as f:
    raw = f.read()

# Split into 3 prodi sections
sections = re.split(r"={20,}\nFILE: (.+?)\n={20,}", raw)
# sections[0]=preamble, then pairs of (filename, content)
chunks = []
for i in range(1, len(sections), 2):
    fname   = sections[i].strip()
    content = sections[i+1]
    if "INFORMATIKA" in fname and "SISTEM" not in fname and "TEKNOLOGI" not in fname:
        prodi = "Informatika"
    elif "SISTEM INFORMASI" in fname:
        prodi = "Sistem Informasi"
    else:
        prodi = "Teknologi Informasi"
    chunks.append((prodi, content))

# Parse courses from a chunk
# Returns: list of dict {code, name, credits, type, semester}
SKIP_CODES = {"LMTR999","MG4191","MG4192","MG4193"}

def parse_courses(prodi, text):
    results = []
    seen    = set()
    semester = 0
    # detect semester headers
    sem_re  = re.compile(r"Semester\s+(\d+)", re.IGNORECASE)
    pre_re  = re.compile(r"pre.requisite", re.IGNORECASE)
    # match course lines: number code name ... credits type
    # pattern: starts with digit space code
    row_re  = re.compile(
        r"^\s*\d+\s+([A-Z][A-Z0-9]+\d+[A-Z0-9]*)\s+(.+?)\s+(\d+)\s+(Pre-requisite|General|Basic|Major|Elective)",
        re.IGNORECASE
    )
    for line in text.splitlines():
        sm = sem_re.search(line)
        if sm:
            semester = int(sm.group(1))
            continue
        if pre_re.search(line) and semester == 0:
            semester = 0
        m = row_re.match(line)
        if not m:
            continue
        code    = m.group(1).strip()
        rawname = m.group(2).strip()
        credits = int(m.group(3))
        ctype   = m.group(4).strip().title().replace("Pre-Requisite","Pre-requisite")
        # clean name: take only the Indonesian/first part before "/"
        # strip trailing ATL markers (single letter preceded by space)
        name = rawname.split("/")[0].strip()
        name = re.sub(r'\s+[a-c]$', '', name).strip()
        if code in SKIP_CODES:
            continue
        key = (code, prodi)
        if key not in seen:
            seen.add(key)
            results.append({"code":code,"name":name,"credits":credits,
                            "type":ctype,"semester":semester,"prodi":prodi})
    return results

all_courses = []
for prodi, text in chunks:
    all_courses.extend(parse_courses(prodi, text))

# Deduplicate codes globally for courses table (one row per code)
global_courses = {}
for c in all_courses:
    if c["code"] not in global_courses:
        global_courses[c["code"]] = c

# ── Parse prerequisites ──────────────────────────────────────
# Format: "- [CODE] Name - N credit(s)"
prereq_re = re.compile(r"-\s+\[([A-Z][A-Z0-9]+\d[A-Z0-9]*)\]")

def parse_prereqs(text):
    prereqs = {}   # course_code -> set of prereq codes
    cur_code = None
    row_re = re.compile(
        r"^\s*\d+\s+([A-Z][A-Z0-9]+\d+[A-Z0-9]*)\s+.+?\s+\d+\s+(Pre-requisite|General|Basic|Major|Elective)",
        re.IGNORECASE
    )
    for line in text.splitlines():
        m = row_re.match(line)
        if m:
            cur_code = m.group(1).strip()
        # find prereq references on any line
        for pm in prereq_re.finditer(line):
            pcode = pm.group(1)
            if cur_code and pcode != cur_code:
                prereqs.setdefault(cur_code, set()).add(pcode)
    return prereqs

all_prereqs = {}  # course_code -> set of prereq_codes
for prodi, text in chunks:
    p = parse_prereqs(text)
    for code, pset in p.items():
        all_prereqs.setdefault(code, set()).update(pset)

# Filter: only keep prereqs where both codes exist in global_courses
filtered_prereqs = {}
for code, pset in all_prereqs.items():
    valid = {p for p in pset if p in global_courses and p != code}
    if valid and code in global_courses:
        filtered_prereqs[code] = valid

# ── Generate SQL ─────────────────────────────────────────────
def q(s):
    return "'" + str(s).replace("'","''") + "'"

lines = []
lines.append("-- ============================================================")
lines.append("-- Seed 003: Lecturers, Courses, Curriculums, Prerequisites")
lines.append("-- Generated automatically from data_kurikulum PDFs")
lines.append("-- ============================================================")
lines.append("SET FOREIGN_KEY_CHECKS = 0;")
lines.append("")

# Lecturers
lines.append("-- ── LECTURERS ─────────────────────────────────────────────")
lines.append("TRUNCATE TABLE lecturers;")
lines.append("INSERT INTO lecturers (id, name) VALUES")
lec_rows = [f"  ({i+1}, {q(name)})" for i, name in enumerate(LECTURERS)]
lines.append(",\n".join(lec_rows) + ";")
lines.append("")

# Courses
lines.append("-- ── COURSES ────────────────────────────────────────────────")
lines.append("TRUNCATE TABLE courses;")
course_vals = []
for code, c in global_courses.items():
    lid = random.randint(1, N_LEC)
    course_vals.append(
        f"  ({q(code)}, {q(c['name'])}, {q('A')}, {q(LECTURERS[lid-1])}, {lid}, {c['credits']}, 2024)"
    )
lines.append("INSERT INTO courses (code, name, class, lecturer_name, lecturer_id, credits, cohort_target) VALUES")
lines.append(",\n".join(course_vals) + ";")
lines.append("")

# Curriculums
lines.append("-- ── CURRICULUMS ─────────────────────────────────────────────")
lines.append("TRUNCATE TABLE curriculums;")
cur_vals = []
for c in all_courses:
    if c["code"] not in global_courses:
        continue
    cur_vals.append(
        f"  ({q(c['code'])}, {q(c['prodi'])}, {c['semester']}, {q(c['type'])})"
    )
lines.append("INSERT INTO curriculums (course_code, study_program, semester, course_type) VALUES")
lines.append(",\n".join(cur_vals) + ";")
lines.append("")

# Prerequisites
lines.append("-- ── PREREQUISITES ───────────────────────────────────────────")
lines.append("TRUNCATE TABLE course_prerequisites;")
pre_vals = []
for code, pset in filtered_prereqs.items():
    for pcode in sorted(pset):
        pre_vals.append(f"  ({q(code)}, {q(pcode)})")
if pre_vals:
    lines.append("INSERT INTO course_prerequisites (course_code, prerequisite_code) VALUES")
    lines.append(",\n".join(pre_vals) + ";")
lines.append("")
lines.append("SET FOREIGN_KEY_CHECKS = 1;")
lines.append("-- Done.")

with open(OUT, "w", encoding="utf-8") as f:
    f.write("\n".join(lines))

print(f"Done. Courses: {len(global_courses)}, Curricula rows: {len(cur_vals)}, Prereqs: {sum(len(v) for v in filtered_prereqs.values())}")
print(f"Output: {OUT}")
