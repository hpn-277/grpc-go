# Super & Salary Sacrifice CRUD - 5 Hour MVP

## Phase 1: Foundation (30 min)
- [x] 1.1 Initialize Go module + folder structure
- [x] 1.2 Docker Compose (PostgreSQL only)
- [x] 1.3 Basic .env and Makefile
- [x] 1.4 Config loader + Database connection

## Phase 2: Shared Domain (20 min)
- [x] 2.1 Email value object
- [x] 2.2 Money value object (cents-based)
- [x] 2.3 Simple UUID helpers

## Phase 3: Database (30 min)
- [x] 3.1 Migration runner (simple version)
- [x] 3.2 Users table migration
- [x] 3.3 Super accounts table migration
- [x] 3.4 Salary sacrifices table migration

## Phase 4: User Feature (60 min)
- [x] 4.1 User domain model (simple struct)
- [x] 4.2 UserRepository interface
- [x] 4.3 GORM UserRepository implementation
- [x] 4.4 UserService (Create, Get, List only)
- [x] 4.5 user.proto (3 endpoints: Create, Get, List)
- [x] 4.6 gRPC handler for UserService

## Phase 5: Super Account Feature (45 min)
- [ ] 5.1 SuperAccount domain model
- [ ] 5.2 SuperAccountRepository interface
- [ ] 5.3 GORM SuperAccountRepository implementation
- [ ] 5.4 SuperAccountService (Create, Get, List by User)
- [ ] 5.5 super_account.proto (3 endpoints)
- [ ] 5.6 gRPC handler for SuperAccountService

## Phase 6: Salary Sacrifice Feature (45 min)
- [ ] 6.1 SalarySacrifice domain model
- [ ] 6.2 SalarySacrificeRepository interface
- [ ] 6.3 GORM SalarySacrificeRepository implementation
- [ ] 6.4 SalarySacrificeService (Create, Get, List by SuperAccount)
- [ ] 6.5 salary_sacrifice.proto (3 endpoints)
- [ ] 6.6 gRPC handler for SalarySacrificeService

## Phase 7: Server & Testing (60 min)
- [ ] 7.1 Main server with dependency injection
- [ ] 7.2 Generate protobuf code
- [ ] 7.3 Wire up all gRPC services
- [ ] 7.4 Test with grpcurl (manual smoke test)
- [ ] 7.5 Basic README with usage examples

---

## Out of Scope (for 5-hour MVP)
- ❌ Authentication/JWT (add later)
- ❌ Password hashing (add later)
- ❌ Update/Delete operations (add later)
- ❌ Domain events (add later)
- ❌ Unit tests (add later)
- ❌ Tax calculator (add later)
- ❌ Middleware/interceptors (add later)
- ❌ Kubernetes/Terraform (add later)

## MVP Scope
✅ **3 entities**: User, SuperAccount, SalarySacrifice
✅ **3 operations per entity**: Create, Get, List
✅ **Clean architecture**: Domain → Application → Infrastructure → gRPC
✅ **Database**: PostgreSQL with migrations
✅ **Working gRPC server** you can test immediately

---

## Architecture Overview

```
internal/
├── features/
│   ├── user-management/
│   │   ├── domain/          # User model, repository interface
│   │   ├── application/     # UserService
│   │   ├── infrastructure/  # GORM repository
│   │   └── entrypoints/     # gRPC handlers
│   ├── super-management/
│   │   ├── domain/
│   │   ├── application/
│   │   ├── infrastructure/
│   │   └── entrypoints/
│   └── sacrifice-management/
│       ├── domain/
│       ├── application/
│       ├── infrastructure/
│       └── entrypoints/
└── shared/
    ├── domain/              # Email, Money value objects
    └── infrastructure/      # Config, Database
```

## Time Estimates
- **Phase 1-3**: 1h 20min (Foundation + Database)
- **Phase 4-6**: 2h 30min (Three CRUD features)
- **Phase 7**: 1h 10min (Server + Testing)
- **Total**: ~5 hours

## Next Steps
Once you're ready to start, we'll build this in order:
1. Foundation → Shared → Database
2. User feature (end-to-end)
3. SuperAccount feature (end-to-end)
4. SalarySacrifice feature (end-to-end)
5. Wire everything together and test
