# Filter CLI

CLI application untuk memfilter data transaksi dari file CSV berdasarkan range waktu.

## Build

```bash
go build -o filter.exe ./cmd/filter
```

Atau untuk Linux/Mac:
```bash
go build -o filter ./cmd/filter
```

## Usage

```bash
filter -d <directory> -s <start> -e <end>
```

### Arguments

| Flag | Description | Format |
|------|-------------|--------|
| `-d` | Path ke directory yang berisi file CSV | Path string |
| `-s` | Start time (inclusive) | RFC3339 |
| `-e` | End time (exclusive) | RFC3339 |

### Contoh

```bash
filter -d ./data -s 2025-06-28T00:00:00+07:00 -e 2025-07-03T00:00:00+07:00
```

## Output

- **Sukses**: Print `successfully filter the data` dan buat file `filtered_result.csv`
- **Error**: Print `unable to filter the data: <detail error>` dengan exit code non-zero

### Lokasi Output File

File hasil filter disimpan di **current working directory** dengan nama `filtered_result.csv`.
Jika file sudah ada, akan di-**overwrite**.

## Format CSV

### Input
- Tidak ada header
- 4 kolom: `TrxNo`, `TrxDate`, `TrxDetail`, `Amount`
- Contoh: `84344,2025-06-28T09:23:55+07:00,transaction: 84344,1863012`

### Output
Format sama persis dengan input (tanpa header).

## Filter Logic

- **Inclusive start**: `TrxDate >= start` akan diinclude
- **Exclusive end**: `TrxDate < end` akan diinclude
- Formula: `start <= TrxDate < end`

## Efisiensi

- File diurutkan berdasarkan numeric prefix (1_report.csv, 2_report.csv, dst)
- Streaming read (line per line, tidak load seluruh file ke memory)
- Early termination: berhenti memproses saat menemukan `TrxDate >= end`

## Run Tests

```bash
go test ./...
```

## Project Structure

```
filter/
├── cmd/filter/main.go      # CLI entrypoint
├── internal/
│   ├── filter/             # Range filter logic
│   ├── io/                 # File I/O (reader, writer)
│   └── parser/             # CSV & time parsing
├── data/                   # Input CSV files (taruh file CSV di sini)
├── testdata/               # Test CSV files untuk unit test
├── go.mod
└── README.md
```
