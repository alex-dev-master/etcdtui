# etcd Client Package

Полнофункциональная обёртка для работы с etcd v3 API.

## Основные возможности

### 1. Базовые операции (CRUD)
- `Get(key)` - получить значение ключа
- `Put(key, value)` - сохранить ключ-значение
- `Delete(key)` - удалить ключ
- `List(prefix)` - получить все ключи с префиксом
- `GetWithRevision(key, revision)` - получить значение на определённой ревизии
- `DeletePrefix(prefix)` - удалить все ключи с префиксом

### 2. Watch (наблюдение за изменениями)
- `Watch(key, callback)` - следить за изменениями ключа
- `WatchPrefix(prefix, callback)` - следить за всеми ключами с префиксом
- `WatchFromRevision(key, revision, callback)` - следить с определённой ревизии

### 3. Lease & TTL
- `PutWithTTL(key, value, ttl)` - сохранить с автоудалением
- `KeepAlive(leaseID)` - поддерживать lease
- `RevokeLease(leaseID)` - отменить lease
- `GetLeaseInfo(leaseID)` - получить информацию о lease
- `ListLeases()` - список всех активных lease

### 4. Транзакции
- `CompareAndSwap(key, oldValue, newValue)` - атомарное обновление
- `CreateIfNotExists(key, value)` - создать только если не существует
- `UpdateIfExists(key, value)` - обновить только если существует

### 5. Распределённые блокировки
- `AcquireLock(key, ttl)` - получить блокировку
- `ReleaseLock(lock)` - освободить блокировку
- `TryLock(key, ttl)` - попытаться получить блокировку
- `IsLocked(key)` - проверить заблокирован ли ключ

### 6. Утилиты
- `GetClusterStatus()` - статус кластера
- `GetKeyCount()` - общее количество ключей
- `GetKeyCountWithPrefix(prefix)` - количество ключей с префиксом
- `BuildTree(keys)` - построить иерархическое дерево
- `CompactHistory(revision)` - сжать историю
- `HealthCheck()` - проверка доступности

### 7. Аутентификация и авторизация
- `CreateUser(username, password)` - создать пользователя
- `DeleteUser(username)` - удалить пользователя
- `ChangePassword(username, password)` - изменить пароль
- `GrantRole(username, role)` - выдать роль
- `RevokeRole(username, role)` - отозвать роль
- `CreateRole(role)` - создать роль
- `GrantPermission(role, key, rangeEnd, permType)` - выдать права
- `EnableAuth()` / `DisableAuth()` - включить/выключить аутентификацию

## Быстрый старт

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/alex-dev-master/etcdtui/internal/client"
)

func main() {
    // Создание клиента с настройками по умолчанию
    cfg := client.DefaultConfig()
    etcdClient, err := client.New(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer etcdClient.Close()

    ctx := context.Background()

    // Сохранить ключ
    err = etcdClient.Put(ctx, "/app/config", "value")
    if err != nil {
        log.Fatal(err)
    }

    // Получить ключ
    kv, err := etcdClient.Get(ctx, "/app/config")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Value: %s\n", kv.Value)
}
```

## Настройка подключения

### Базовое подключение
```go
cfg := &client.Config{
    Endpoints: []string{"localhost:2379"},
    DialTimeout: 5 * time.Second,
    RequestTimeout: 10 * time.Second,
}
```

### С аутентификацией
```go
cfg := &client.Config{
    Endpoints: []string{"etcd.example.com:2379"},
    Username: "admin",
    Password: "secret",
    DialTimeout: 5 * time.Second,
}
```

### С TLS
```go
cfg := &client.Config{
    Endpoints: []string{"https://etcd.example.com:2379"},
    TLS: &client.TLSConfig{
        Enabled: true,
        CertFile: "/path/to/client.crt",
        KeyFile: "/path/to/client.key",
        CAFile: "/path/to/ca.crt",
    },
}
```

## Примеры использования

Подробные примеры смотрите в файле `examples.go`.

### Watch для изменений
```go
err := etcdClient.WatchPrefix(ctx, "/config/", func(event *client.WatchEvent) {
    if event.Type == client.EventTypePut {
        fmt.Printf("Updated: %s = %s\n", event.Key, event.Value)
    } else {
        fmt.Printf("Deleted: %s\n", event.Key)
    }
})
```

### Distributed Lock
```go
lock, err := etcdClient.AcquireLock(ctx, "/locks/resource", 30*time.Second)
if err != nil {
    return err
}
defer etcdClient.ReleaseLock(ctx, lock)

// Критическая секция
// ...
```

### TTL ключи
```go
lease, err := etcdClient.PutWithTTL(ctx, "/session/user1", "active", 60*time.Second)
if err != nil {
    return err
}

// Поддерживать сессию
keepAliveCh, _ := etcdClient.KeepAlive(ctx, lease.ID)
for range keepAliveCh {
    // Lease продлён
}
```

## Обработка ошибок

Все методы возвращают ошибки в формате `error`. Используйте стандартные методы Go для обработки:

```go
kv, err := etcdClient.Get(ctx, "/key")
if err != nil {
    // Обработка ошибки
    log.Printf("Failed to get key: %v", err)
    return err
}
```

## Context и таймауты

Все методы принимают `context.Context`. Клиент автоматически применяет таймауты из конфигурации, но вы можете использовать собственные:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

err := etcdClient.Put(ctx, "/key", "value")
```

## Thread Safety

Клиент потокобезопасен и может использоваться из нескольких горутин одновременно.

## Производительность

- Используйте `List()` вместо множественных `Get()` для получения нескольких ключей
- Для массового удаления используйте `DeletePrefix()` вместо множественных `Delete()`
- Watch операции выполняются в фоне и не блокируют
- Lease keep-alive работает автоматически в фоновом режиме
