# Nasabah Kredit API

## Endpoints

---

- **endpoint**: **/api/nasabah-kredit**
  **method**: **GET**
  **description**: List nasabah kredit 
  **query parameters (untuk pencarian)**:
   ```json
   {
     "cif": "string",
     "nama": "string",
     "nik": "string",
     "nomor_rekening": "string",
   }
   ```
  **response**:
   ```json
   {
     "data": [
       {
         "cif": "string",
         "nik": "string",
         "nama": "string",
         "tempat_lahir": "string",
         "tanggal_lahir": "integer (unix ms)",
         "jenis_kelamin": "string",
         "alamat": "string",
         "phone": "string",
         "email": "string",
         "jenis_pinjaman" : "string",
         "nomor_rekening": "string",
         "unit": "string",
         "kolektibilitas": "string",
         "plafond": "number",
         "sisa_tagihan": "number", // total pokok + bunga yang belum dibayar
         "tunggakan": "number", // jumlah angsuran yang belum dibayar
         "created_at": "integer (unix ms)",
         "updated_at": "integer (unix ms)"
       }
     ],
   }
   ```

---

- endpoint: /api/nasabah-kredit
  method: POST
  description: Create a new nasabah kredit.
  request body:
   ```json
   {
     "nik": "string (required)",
     "nama": "string (required)",
     "tempat_lahir": "string (required)",
     "tanggal_lahir": "string (required, format: YYYY-MM-DD)",
     "jenis_kelamin": "string (required, enum: LAKI_LAKI | PEREMPUAN)",
     "alamat": "string (required)",
     "phone": "string (required)",
     "email": "string",
     "nomor_rekening": "string",
     "pekerjaan": "string",
     "penghasilan_perbulan": "number"
   }
   ```
  response:
   ```json
   {
     "data": {
       "id": "string",
       "nik": "string",
       "nama": "string",
       "tempat_lahir": "string",
       "tanggal_lahir": "integer (unix ms)",
       "jenis_kelamin": "string",
       "alamat": "string",
       "phone": "string",
       "email": "string",
       "pekerjaan": "string",
       "penghasilan_perbulan": "number",
       "status": "string",
       "kolektibilitas": "string",
       "created_at": "integer (unix ms)",
       "updated_at": "integer (unix ms)"
     },
     "paging": null,
     "errors": null
   }
   ```

---

- endpoint: /api/nasabah-kredit/{id}
  method: GET
  description: Get nasabah kredit detail by ID.
  path parameters:
   ```json
   {
     "id": "string (required)"
   }
   ```
  response:
   ```json
   {
     "data": {
       "id": "string",
       "nik": "string",
       "nama": "string",
       "tempat_lahir": "string",
       "tanggal_lahir": "integer (unix ms)",
       "jenis_kelamin": "string",
       "alamat": "string",
       "phone": "string",
       "email": "string",
       "pekerjaan": "string",
       "penghasilan_perbulan": "number",
       "status": "string",
       "kredit": [
         {
           "id": "string",
           "jumlah_pinjaman": "number",
           "tenor_bulan": "integer",
           "bunga_persen": "number",
           "angsuran_perbulan": "number",
           "tanggal_pencairan": "integer (unix ms)",
           "tanggal_jatuh_tempo": "integer (unix ms)",
           "status_kredit": "string"
         }
       ],
       "created_at": "integer (unix ms)",
       "updated_at": "integer (unix ms)"
     },
     "paging": null,
     "errors": null
   }
   ```

---

- endpoint: /api/nasabah-kredit/{id}
  method: PUT
  description: Update nasabah kredit data by ID.
  path parameters:
   ```json
   {
     "id": "string (required)"
   }
   ```
  request body:
   ```json
   {
     "nama": "string",
     "tempat_lahir": "string",
     "tanggal_lahir": "string (format: YYYY-MM-DD)",
     "jenis_kelamin": "string (enum: LAKI_LAKI | PEREMPUAN)",
     "alamat": "string",
     "phone": "string",
     "email": "string",
     "pekerjaan": "string",
     "penghasilan_perbulan": "number",
     "status": "string (enum: AKTIF | TIDAK_AKTIF | BLACKLIST)"
   }
   ```
  response:
   ```json
   {
     "data": {
       "id": "string",
       "nik": "string",
       "nama": "string",
       "tempat_lahir": "string",
       "tanggal_lahir": "integer (unix ms)",
       "jenis_kelamin": "string",
       "alamat": "string",
       "phone": "string",
       "email": "string",
       "pekerjaan": "string",
       "penghasilan_perbulan": "number",
       "status": "string",
       "created_at": "integer (unix ms)",
       "updated_at": "integer (unix ms)"
     },
     "paging": null,
     "errors": null
   }
   ```

---

- endpoint: /api/nasabah-kredit/{id}
  method: DELETE
  description: Delete nasabah kredit by ID.
  path parameters:
   ```json
   {
     "id": "string (required)"
   }
   ```
  response:
   ```json
   {
     "data": null,
     "paging": null,
     "errors": null
   }
   ```

---

- endpoint: /api/nasabah-kredit/search
  method: GET
  description: Search nasabah kredit based on query parameters.
  query parameters:
   ```json
   {
     "query": "string",
     "status": "string",
     "limit": "integer"
   }
   ```
  response:
   ```json
   {
     "data": [
       {
         "id": "string",
         "nik": "string",
         "nama": "string",
         "phone": "string",
         "status": "string"
       }
     ],
     "paging": {
       "page": "integer",
       "size": "integer",
       "total_item": "integer",
       "total_page": "integer"
     },
     "errors": null
   }
   ```

---

- endpoint: /api/nasabah-kredit/{id}/kredit
  method: GET
  description: List all kredit (loan) records for a specific nasabah.
  path parameters:
   ```json
   {
     "id": "string (required)"
   }
   ```
  query parameters:
   ```json
   {
     "status_kredit": "string",
     "page": "integer",
     "size": "integer"
   }
   ```
  response:
   ```json
   {
     "data": [
       {
         "id": "string",
         "nasabah_id": "string",
         "jumlah_pinjaman": "number",
         "tenor_bulan": "integer",
         "bunga_persen": "number",
         "angsuran_perbulan": "number",
         "tanggal_pencairan": "integer (unix ms)",
         "tanggal_jatuh_tempo": "integer (unix ms)",
         "status_kredit": "string (enum: AKTIF | LUNAS | MACET | DITOLAK)"
       }
     ],
     "paging": {
       "page": "integer",
       "size": "integer",
       "total_item": "integer",
       "total_page": "integer"
     },
     "errors": null
   }
   ```

---

- endpoint: /api/nasabah-kredit/{id}/kredit
  method: POST
  description: Create a new kredit (loan) for a nasabah.
  path parameters:
   ```json
   {
     "id": "string (required)"
   }
   ```
  request body:
   ```json
   {
     "jumlah_pinjaman": "number (required, min: 1)",
     "tenor_bulan": "integer (required, min: 1)",
     "bunga_persen": "number (required, min: 0)",
     "tanggal_pencairan": "string (required, format: YYYY-MM-DD)",
     "tujuan_pinjaman": "string"
   }
   ```
  response:
   ```json
   {
     "data": {
       "id": "string",
       "nasabah_id": "string",
       "jumlah_pinjaman": "number",
       "tenor_bulan": "integer",
       "bunga_persen": "number",
       "angsuran_perbulan": "number",
       "tanggal_pencairan": "integer (unix ms)",
       "tanggal_jatuh_tempo": "integer (unix ms)",
       "tujuan_pinjaman": "string",
       "status_kredit": "string"
     },
     "paging": null,
     "errors": null
   }
   ```

---

## Error Responses

All endpoints return errors in this format:

```json
{
  "success": false,
  "message": "string"
}
```

| HTTP Status | Description              |
|-------------|--------------------------|
| 400         | Bad request / validation |
| 401         | Unauthorized             |
| 403         | Forbidden                |
| 404         | Not found                |
| 422         | Unprocessable entity     |
| 500         | Internal server error    |
