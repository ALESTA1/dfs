# Distributed File System

A robust distributed file system implementation inspired by a CMU course project. The system features a centralized Naming Server architecture with distributed Storage Servers, supporting file replication and lock-based consistency.

## 🏗️ Architecture

### Naming Server
- Acts as the central coordinator between clients and storage servers
- Manages file locations and lock distribution
- Single instance running at any time
- Handles server registration and file directory updates

### Storage Servers
- Multiple instances can run simultaneously
- Register with Naming Server on startup
- Maintain local file directories
- Support file replication and consistency mechanisms

## 🔐 Lock Management

The system implements two types of locks for file access:
- **Shared Lock**: Required for read operations
- **Exclusive Lock**: Required for write operations

> **Note**: Due to architectural constraints, shared locks are treated as write operations, and exclusive locks are treated as read operations at the implementation level.

## 📦 Features

### Automatic Replication
- New replicas are created after every 20 shared locks on a file
- Replicas are distributed across available servers
- Ensures non-blocking access for multiple clients

### Consistency Management
- Exclusive lock requests trigger replica consolidation
- Maintains single source of truth for write operations
- Automatic conflict resolution during server registration

## 🚀 Getting Started

### Prerequisites
- Go runtime environment
- Available network ports for services
- Storage space for file systems

### Starting the Naming Server

Navigate to the `naming` directory and run:

```bash
go run . <Client_Service_Port> <Slave_Registration_Port>
```

Example:
```bash
go run . 8080 8081
```

### Starting a Storage Server

Navigate to the `storage` directory and run:

```bash
go run . <Client_Service_Port> <Command_Port> <Registration_Port> <Storage_Directory>
```

Example:
```bash
go run . 9090 9091 8081 storage1
```

## 💡 Implementation Notes

### Storage Directory Handling
- Automatically created if it doesn't exist
- Preserves existing contents if directory is present
- Handles file conflicts during registration

### Client Expectations
1. Clients must request appropriate locks before file operations
2. Locks should be released after operations complete
3. Cooperative access model is assumed

## 🔍 System Assumptions

1. Clients follow the locking protocol:
   - Acquire appropriate locks before operations
   - Release locks after completing operations
   - Use shared locks for reads and exclusive locks for writes

2. File Replication:
   - System maintains consistency across replicas
   - Latest version access is guaranteed
   - Non-blocking operations during normal usage

## 📚 Technical Details

### Replication Mechanism
- Triggered after 20 shared locks
- Distributed across available servers
- Automatic cleanup during exclusive access

### Consistency Protocol
- Single-writer, multiple-reader model
- Automatic conflict resolution
- Replica consolidation during writes

## 🤝 Contributing

Feel free to submit issues and enhancement requests.

## 🙏 Acknowledgments

This project is inspired by a CMU course project. For detailed project requirements, refer to the original course materials.
