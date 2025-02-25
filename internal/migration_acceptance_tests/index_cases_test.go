package migration_acceptance_tests

import (
	"github.com/stripe/pg-schema-diff/pkg/diff"
)

var indexAcceptanceTestCases = []acceptanceTestCase{
	{
		name: "No-op",
		oldSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT PRIMARY KEY,
				foo VARCHAR(255),
				bar TEXT,
				fizz INT
			);
			CREATE INDEX some_idx ON foobar USING hash (foo);
			CREATE UNIQUE INDEX some_other_idx ON foobar (bar DESC, fizz);
			`,
		},
		newSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT PRIMARY KEY,
				foo VARCHAR(255),
				bar TEXT,
				fizz INT
			);
			CREATE INDEX some_idx ON foobar USING hash (foo);
			CREATE UNIQUE INDEX some_other_idx ON foobar (bar DESC, fizz);
			`,
		},
	},
	{
		name: "Add a normal index",
		oldSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT PRIMARY KEY,
				foo VARCHAR(255)
			);
			`,
		},
		newSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT PRIMARY KEY,
				foo VARCHAR(255)
			);
			CREATE INDEX some_idx ON foobar(id DESC, foo);
			`,
		},
		expectedHazardTypes: []diff.MigrationHazardType{
			diff.MigrationHazardTypeIndexBuild,
		},
	},
	{
		name: "Add a hash index",
		oldSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT PRIMARY KEY,
				foo VARCHAR(255)
			);
			`,
		},
		newSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT PRIMARY KEY,
				foo VARCHAR(255)
			);
			CREATE INDEX some_idx ON foobar USING hash (id);
			`,
		},
		expectedHazardTypes: []diff.MigrationHazardType{
			diff.MigrationHazardTypeIndexBuild,
		},
	},
	{
		name: "Add a normal index with quoted names",
		oldSchemaDDL: []string{
			`
			CREATE TABLE "Foobar"(
			    id INT PRIMARY KEY,
				"Foo" VARCHAR(255)
			);
			`,
		},
		newSchemaDDL: []string{
			`
			CREATE TABLE "Foobar"(
			    id INT PRIMARY KEY,
				"Foo" VARCHAR(255)
			);
			CREATE INDEX "Some_idx" ON "Foobar"(id, "Foo");
			`,
		},
		expectedHazardTypes: []diff.MigrationHazardType{
			diff.MigrationHazardTypeIndexBuild,
		},
	},
	{
		name: "Add a unique index",
		oldSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT PRIMARY KEY,
			    foo VARCHAR(255) NOT NULL
			);
			`,
		},
		newSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT PRIMARY KEY,
			    foo VARCHAR(255) NOT NULL
			);
			CREATE UNIQUE INDEX some_unique_idx ON foobar(foo);
			`,
		},
		expectedHazardTypes: []diff.MigrationHazardType{
			diff.MigrationHazardTypeIndexBuild,
		},
	},
	{
		name: "Add a primary key",
		oldSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT
			);
			`,
		},
		newSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT PRIMARY KEY
			);
			`,
		},
		expectedHazardTypes: []diff.MigrationHazardType{
			diff.MigrationHazardTypeAcquiresAccessExclusiveLock,
			diff.MigrationHazardTypeIndexBuild,
		},
	},
	{
		name: "Add a primary key on NOT NULL column",
		oldSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT NOT NULL
			);
			`,
		},
		newSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT NOT NULL PRIMARY KEY
			);
			`,
		},
		expectedHazardTypes: []diff.MigrationHazardType{
			diff.MigrationHazardTypeIndexBuild,
		},
	},
	{
		name: "Add a primary key when the index already exists",
		oldSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT
			);
			CREATE UNIQUE INDEX foobar_primary_key ON foobar(id);
			`,
		},
		newSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT
			);
			CREATE UNIQUE INDEX foobar_primary_key ON foobar(id);
			ALTER TABLE foobar ADD CONSTRAINT foobar_primary_key PRIMARY KEY USING INDEX foobar_primary_key;
			`,
		},
		expectedHazardTypes: []diff.MigrationHazardType{
			diff.MigrationHazardTypeAcquiresAccessExclusiveLock,
		},
	},
	{
		name: "Add a primary key when the index already exists but has a name different to the constraint",
		oldSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT
			);
			CREATE UNIQUE INDEX foobar_idx ON foobar(id);
			`,
		},
		newSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT
			);
			CREATE UNIQUE INDEX foobar_idx ON foobar(id);
			-- This renames the index
			ALTER TABLE foobar ADD CONSTRAINT foobar_primary_key PRIMARY KEY USING INDEX foobar_idx;
			`,
		},
		expectedHazardTypes: []diff.MigrationHazardType{
			diff.MigrationHazardTypeAcquiresAccessExclusiveLock,
			diff.MigrationHazardTypeIndexDropped,
			diff.MigrationHazardTypeIndexBuild,
		},
	},
	{
		name: "Add a primary key when the index already exists but is not unique",
		oldSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT
			);
			CREATE INDEX foobar_idx ON foobar(id);
			`,
		},
		newSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT
			);
			CREATE UNIQUE INDEX foobar_primary_key ON foobar(id);
			ALTER TABLE foobar ADD CONSTRAINT foobar_primary_key PRIMARY KEY USING INDEX foobar_primary_key;
			`,
		},
		expectedHazardTypes: []diff.MigrationHazardType{
			diff.MigrationHazardTypeAcquiresAccessExclusiveLock,
			diff.MigrationHazardTypeIndexDropped,
			diff.MigrationHazardTypeIndexBuild,
		},
	},
	{
		name: "Delete a normal index",
		oldSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT PRIMARY KEY,
				foo VARCHAR(255) NOT NULL
			);
			CREATE INDEX some_inx ON foobar(id, foo);
			`,
		},
		newSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT PRIMARY KEY,
				foo VARCHAR(255)
			);
			`,
		},
		expectedHazardTypes: []diff.MigrationHazardType{
			diff.MigrationHazardTypeIndexDropped,
		},
	},
	{
		name: "Delete a normal index with quoted names",
		oldSchemaDDL: []string{
			`
			CREATE TABLE "Foobar"(
			    id INT PRIMARY KEY,
				"Foo" VARCHAR(255)
			);
			CREATE INDEX "Some_idx" ON "Foobar"(id, "Foo");
			`,
		},
		newSchemaDDL: []string{
			`
			CREATE TABLE "Foobar"(
			    id INT PRIMARY KEY,
				"Foo" VARCHAR(255) NOT NULL
			);
			`,
		},
		expectedHazardTypes: []diff.MigrationHazardType{
			diff.MigrationHazardTypeAcquiresAccessExclusiveLock,
			diff.MigrationHazardTypeIndexDropped,
		},
	},
	{
		name: "Delete a unique index",
		oldSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT PRIMARY KEY,
			    foo VARCHAR(255) NOT NULL
			);
			CREATE UNIQUE INDEX some_unique_idx ON foobar(foo);
			`,
		},
		newSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT PRIMARY KEY,
			    foo VARCHAR(255) NOT NULL
			);
			`,
		},
		expectedHazardTypes: []diff.MigrationHazardType{
			diff.MigrationHazardTypeIndexDropped,
		},
	},
	{
		name: "Delete a primary key",
		oldSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT PRIMARY KEY
			);
			`,
		},
		newSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT
			);
			`,
		},
		expectedHazardTypes: []diff.MigrationHazardType{
			diff.MigrationHazardTypeAcquiresAccessExclusiveLock,
			diff.MigrationHazardTypeIndexDropped,
		},
	},
	{
		name: "Change an index (with a really long name) columns",
		oldSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT PRIMARY KEY,
				foo TEXT NOT NULL,
			    bar BIGINT NOT NULL
			);
			CREATE UNIQUE INDEX some_idx_with_a_really_long_name_that_is_nearly_61_chars ON foobar(foo, bar)
			`,
		},
		newSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT,
				foo TEXT NOT NULL,
			    bar BIGINT NOT NULL
			);
			CREATE UNIQUE INDEX some_idx_with_a_really_long_name_that_is_nearly_61_chars ON foobar(foo)
			`,
		},
		expectedHazardTypes: []diff.MigrationHazardType{
			diff.MigrationHazardTypeAcquiresAccessExclusiveLock,
			diff.MigrationHazardTypeIndexDropped,
			diff.MigrationHazardTypeIndexBuild,
		},
	},
	{
		name: "Change an index type",
		oldSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT PRIMARY KEY,
				foo TEXT NOT NULL,
			    bar BIGINT NOT NULL
			);
			CREATE INDEX some_idx ON foobar (foo)
			`,
		},
		newSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT,
				foo TEXT NOT NULL,
			    bar BIGINT NOT NULL
			);
			CREATE INDEX some_idx ON foobar USING hash (foo)
			`,
		},
		expectedHazardTypes: []diff.MigrationHazardType{
			diff.MigrationHazardTypeAcquiresAccessExclusiveLock,
			diff.MigrationHazardTypeIndexDropped,
			diff.MigrationHazardTypeIndexBuild,
		},
	},
	{
		name: "Change an index column ordering",
		oldSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT PRIMARY KEY,
				foo TEXT NOT NULL,
			    bar BIGINT NOT NULL
			);
			CREATE INDEX some_idx ON foobar (foo, bar)
			`,
		},
		newSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT,
				foo TEXT NOT NULL,
			    bar BIGINT NOT NULL
			);
			CREATE INDEX some_idx ON foobar (foo DESC, bar)
			`,
		},
		expectedHazardTypes: []diff.MigrationHazardType{
			diff.MigrationHazardTypeAcquiresAccessExclusiveLock,
			diff.MigrationHazardTypeIndexDropped,
			diff.MigrationHazardTypeIndexBuild,
		},
	},
	{
		name: "Delete columns and associated index",
		oldSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT PRIMARY KEY,
				foo TEXT NOT NULL,
			    bar BIGINT NOT NULL
			);
			CREATE UNIQUE INDEX some_idx ON foobar(foo, bar);
			`,
		},
		newSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT
			);
			`,
		},
		expectedHazardTypes: []diff.MigrationHazardType{
			diff.MigrationHazardTypeAcquiresAccessExclusiveLock,
			diff.MigrationHazardTypeDeletesData,
			diff.MigrationHazardTypeIndexDropped,
		},
	},
	{
		name: "Switch primary key and make old key nullable",
		oldSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT NOT NULL PRIMARY KEY,
				foo TEXT NOT NULL
			);
			CREATE UNIQUE INDEX some_idx ON foobar(foo);
			`,
		},
		newSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT,
			    foo INT NOT NULL PRIMARY KEY
			);
			`,
		},
		expectedHazardTypes: []diff.MigrationHazardType{
			diff.MigrationHazardTypeAcquiresAccessExclusiveLock,
			diff.MigrationHazardTypeIndexDropped,
			diff.MigrationHazardTypeIndexBuild,
			diff.MigrationHazardTypeImpactsDatabasePerformance,
		},
	},
	{
		name: "Switch primary key with quoted name",
		oldSchemaDDL: []string{
			`
			CREATE TABLE "Foobar"(
			    "Id" INT NOT NULL PRIMARY KEY,
				foo TEXT NOT NULL
			);
			`,
		},
		newSchemaDDL: []string{
			`
			CREATE TABLE "Foobar"(
			    "Id" INT,
			    foo INT NOT NULL PRIMARY KEY
			);
			`,
		},
		expectedHazardTypes: []diff.MigrationHazardType{
			diff.MigrationHazardTypeAcquiresAccessExclusiveLock,
			diff.MigrationHazardTypeIndexDropped,
			diff.MigrationHazardTypeIndexBuild,
			diff.MigrationHazardTypeImpactsDatabasePerformance,
		},
	},
	{
		name: "Switch primary key when the original primary key constraint has a non-default name",
		oldSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT NOT NULL,
				foo TEXT NOT NULL
			);
			CREATE UNIQUE INDEX unique_idx ON foobar(id);
			ALTER TABLE foobar ADD CONSTRAINT non_default_primary_key PRIMARY KEY USING INDEX unique_idx;
			`,
		},
		newSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT,
			    foo INT NOT NULL PRIMARY KEY
			);
			`,
		},
		expectedHazardTypes: []diff.MigrationHazardType{
			diff.MigrationHazardTypeAcquiresAccessExclusiveLock,
			diff.MigrationHazardTypeImpactsDatabasePerformance,
			diff.MigrationHazardTypeIndexDropped,
			diff.MigrationHazardTypeIndexBuild,
		},
	},
	{
		name: "Alter primary key columns (name stays same)",
		oldSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT NOT NULL,
				foo TEXT NOT NULL
			);
			CREATE UNIQUE INDEX unique_idx ON foobar(id);
			ALTER TABLE foobar ADD CONSTRAINT non_default_primary_key PRIMARY KEY USING INDEX unique_idx;
			`,
		},
		newSchemaDDL: []string{
			`
			CREATE TABLE foobar(
			    id INT,
			    foo INT NOT NULL
			);
			CREATE UNIQUE INDEX unique_idx ON foobar(id, foo);
			ALTER TABLE foobar ADD CONSTRAINT non_default_primary_key PRIMARY KEY USING INDEX unique_idx;
			`,
		},
		expectedHazardTypes: []diff.MigrationHazardType{
			diff.MigrationHazardTypeAcquiresAccessExclusiveLock,
			diff.MigrationHazardTypeImpactsDatabasePerformance,
			diff.MigrationHazardTypeIndexDropped,
			diff.MigrationHazardTypeIndexBuild,
		},
	},
}

func (suite *acceptanceTestSuite) TestIndexAcceptanceTestCases() {
	suite.runTestCases(indexAcceptanceTestCases)
}
