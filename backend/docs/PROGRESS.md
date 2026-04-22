# Project Progress Form: Sistem Kontrak Mata Kuliah

*Changelog and tracking of tasks completed.*

## Version: Domain & Handlers Skeleton (Step 2)

### Completed Features:

**1. Defined Core Domain Models (`domain` layer):**
- **`domain/user.go`**: 
  - Added `User` struct including fields for `NIM`, `Name`, `Faculty`, `StudyProgram`, `CohortYear` (Angkatan), and `Role`, nicely mapped with `json` tags. 
  - Implemented initial skeletons for `UserRepository` and `UserUsecase` interfaces.
- **`domain/course.go`**: 
  - Added `Course` bounding structure parsing properties such as `Code`, `Name`, `Class`, `Credits`, etc.
  - Linked `Schedule` structure with timing properties like `Day`, `StartTime`, `EndTime`, `Room`, and `AdditionalNotes`. 
  - Implemented skeleton `CourseRepository` and `CourseUsecase` interfaces.

**2. Refactored HTTP Handlers (`delivery/http` layer):**
- Standardized declarative handler definitions following Clean Architecture via interface injection structures.
- Migrated out generic string logic from generic loops into decoupled method pointers mapping specifically to:
  - `AdminHandler`
  - `StudentHandler`
- Preserved strict `net/http` implementations while exposing pointer references (`*http.Request`). 

**3. Wired ServeMux Setup (`main.go`):**
- Instantiated decoupled application controllers using basic DI mapping passing explicit `nil` stubs for internal business rules.
- Transitioned path bindings integrating handler methods properly against Role Middleware implementations using Go 1.22's standardized `.HandleFunc("METHOD /path", handler)` syntax structure.

### Next Steps / Pending Elements:

- Prepare specific database connections and DTO models.
- Implement strictly defined concrete business rule functions across the `usecase` directories binding correctly toward generic requests mapping models.
- Set up data layers and implementations connecting directly via standard SQL logic (database mappings) located in `repository`.
