package model

var CreatePersonTableSql = `
create table Person (
    id string primary key not null,
    createdOn integer not null,
    modifiedOn integer not null,
    publicKey string not null,
    name string not null,
    email string not null,
    isMajor bool not null,
    isSigned bool not null,
    saldo integer not null,
    biographyHash string not null,
    biographyFormat string not null,
    institution string not null,
    telephone string not null,
    address string not null,
    zipCode string not null,
    country string not null,
    governmentId string not null
)`
