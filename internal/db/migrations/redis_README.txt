# Redis setup for Hauling App
# Geospatial driver locations will be stored in the 'drivers' key using GEOADD/GEOSEARCH
# Example:
# GEOADD drivers <longitude> <latitude> <driver_id>

# Password reset tokens (if implemented) can be stored as:
# SET reset:<user_id> <token> EX <expiry_seconds>

# No migration file is needed for Redis, but this README documents the usage.

# For local development, you may want to flushdb before starting:
# FLUSHDB
