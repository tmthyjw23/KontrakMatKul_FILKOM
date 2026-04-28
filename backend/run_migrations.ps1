# Script untuk menjalankan semua migrations ke MySQL database
$DB_HOST = "127.0.0.1"
$DB_PORT = "3306"
$DB_USER = "root"
$DB_PASSWORD = ""
$DB_NAME = "sistem_kontrak"

# Create database if not exists
$createDbCmd = "mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -e `"CREATE DATABASE IF NOT EXISTS \`$DB_NAME\`;`""
Write-Host "Creating database..."
Invoke-Expression $createDbCmd

# Run all migration files
$migrationFiles = @(
    "000001_create_core_tables.up.sql",
    "000002_add_role_and_password_to_users.up.sql",
    "000003_add_lecturer_and_room_to_courses_and_schedules.up.sql",
    "000003_create_system_settings.up.sql",
    "000004_create_passed_courses.up.sql",
    "000005_create_course_prerequisites.up.sql",
    "000006_add_period_dates.up.sql",
    "seed_users.sql"
)

foreach ($file in $migrationFiles) {
    Write-Host "Running migration: $file"
    $fullPath = "migrations\$file"
    if (Test-Path $fullPath) {
        $cmd = "mysql -h $DB_HOST -P $DB_PORT -u $DB_USER $DB_NAME < `"$fullPath`""
        Invoke-Expression $cmd
        Write-Host "✓ Completed: $file"
    } else {
        Write-Host "✗ File not found: $fullPath"
    }
}

Write-Host "All migrations completed!"
