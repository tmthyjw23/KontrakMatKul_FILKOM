# Frontend Manual

Dokumentasi ini khusus untuk aplikasi frontend di folder `frontend/`.

## Stack

- Next.js App Router
- React
- Tailwind CSS
- Framer Motion
- Zustand
- Axios
- TanStack React Query
- Sonner

## Tujuan Frontend

Frontend berfungsi sebagai dashboard kontrak mata kuliah dengan dua pane:

- kiri: daftar curriculum
- kanan: visual weekly schedule, SKS counter, dan tombol confirm enrollment

## Struktur Utama

```text
frontend/
├─ app/
├─ components/
├─ lib/
├─ src/types/
├─ public/
└─ MANUAL.md
```

Folder penting:

- `app/`: layout, page, provider
- `components/dashboard/`: komponen domain dashboard
- `components/ui/`: komponen UI reusable
- `lib/api/`: axios client
- `lib/hooks/`: React Query hooks
- `lib/store/`: Zustand store
- `lib/utils/`: util jadwal
- `src/types/`: tipe data terpusat

## Cara Menjalankan

```powershell
cd frontend
npm install
npm run dev
```

Default URL:

```text
http://localhost:3000
```

## Environment

Frontend membaca API URL dari:

```text
NEXT_PUBLIC_API_URL
```

Jika tidak diisi, default:

```text
http://localhost:8080/api/v1
```

Contoh `.env.local`:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

## Alur Data

### Curriculum

Curriculum diambil dari backend dengan hook:

- `lib/hooks/useCourses.ts`

Endpoint:

```text
GET /courses
```

React Query akan:

- cache hasil fetch
- mengurangi refetch yang tidak perlu
- memberi state `loading`, `error`, dan `success`

### Enrollment

Saat user menekan tombol confirm:

- frontend membaca `selectedCourses` dari Zustand
- frontend mengirim satu payload enrollment per course ke backend
- jika sukses, selection di-clear
- curriculum di-invalidate dan dapat di-refetch

Hook terkait:

- `lib/hooks/useEnrollment.ts`

## State Management

Zustand digunakan untuk:

- `courses`
- `selectedCourses`
- `totalSks`
- `maxSks`
- visual conflict detection

Store utama:

- `lib/store/useContractStore.ts`

## Tipe Data

Model utama ada di:

- `src/types/course.ts`

Ini penting agar struktur data frontend tetap konsisten dengan backend.

## Login dan Token

Axios akan otomatis mencoba membaca token dari:

- `localStorage.getItem("token")`
- cookie `token`

Token dipasang ke header:

```text
Authorization: Bearer <token>
```

## Komponen Penting

- `app/page.tsx`: dashboard utama
- `components/dashboard/course-card.tsx`: card curriculum
- `components/dashboard/schedule-grid.tsx`: visual planner
- `components/dashboard/sks-counter.tsx`: gauge SKS
- `components/ui/glass-card.tsx`: wrapper glassmorphism

## Validasi Visual di Frontend

Frontend memiliki deteksi visual conflict untuk menandai bentrok jadwal sebelum request dikirim. Ini hanya bantuan UI.

Validasi final tetap terjadi di backend.

## Troubleshooting

### Gagal memuat curriculum

Periksa:

- backend berjalan
- `NEXT_PUBLIC_API_URL` benar
- `GET /api/v1/courses` mengembalikan `200`

### Confirm Enrollment gagal

Periksa:

- token JWT tersedia
- login berhasil
- backend `POST /api/v1/enrollments` aktif

### Cek kualitas code

```powershell
npm run lint
npx tsc --noEmit
```
