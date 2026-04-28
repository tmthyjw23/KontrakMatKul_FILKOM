import pdfplumber
import os

data_dir = r'c:\Users\Dave Jordy\OneDrive\Documents\KULIAHH\SEMESTER 4\KontrakMatKul_FILKOM\backend\docs\data_kurikulum'
out_file = r'c:\Users\Dave Jordy\OneDrive\Documents\KULIAHH\SEMESTER 4\KontrakMatKul_FILKOM\backend\docs\extracted_all.txt'

target_pdfs = [
    'Kurikulum FILKOM ALL PRODI - INFORMATIKA.pdf',
    'Kurikulum FILKOM ALL PRODI - SISTEM INFORMASI 2020.pdf',
    'Kurikulum FILKOM ALL PRODI - TEKNOLOGI INFORMASI.pdf',
]

sep = '=' * 80

with open(out_file, 'w', encoding='utf-8') as out:
    for pdf_name in target_pdfs:
        pdf_path = os.path.join(data_dir, pdf_name)
        out.write('\n\n' + sep + '\n')
        out.write('FILE: ' + pdf_name + '\n')
        out.write(sep + '\n')
        with pdfplumber.open(pdf_path) as pdf:
            out.write('Total pages: ' + str(len(pdf.pages)) + '\n\n')
            for i, page in enumerate(pdf.pages):
                text = page.extract_text()
                out.write('--- PAGE ' + str(i + 1) + ' ---\n')
                if text:
                    out.write(text)
                else:
                    out.write('[No text]')
                out.write('\n\n')

print('Done. Output written to: ' + out_file)
