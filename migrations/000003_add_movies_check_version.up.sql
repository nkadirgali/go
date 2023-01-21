-- make sure that the version to be less than 10
ALTER TABLE movies 
        ADD CONSTRAINT movies_version_check 
        CHECK (version < 10);