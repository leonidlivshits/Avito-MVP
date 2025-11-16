-- 0001_init.up.sql
BEGIN;

-- users
CREATE TABLE IF NOT EXISTS users (
  user_id TEXT PRIMARY KEY,
  username TEXT NOT NULL,
  team_name TEXT NOT NULL,
  skill_level TEXT NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_users_team ON users(team_name);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);

-- teams
CREATE TABLE IF NOT EXISTS teams (
  team_name TEXT PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- team members (optional explicit relation)
CREATE TABLE IF NOT EXISTS team_members (
  team_name TEXT NOT NULL REFERENCES teams(team_name) ON DELETE CASCADE,
  user_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  PRIMARY KEY (team_name, user_id)
);

-- pull requests
CREATE TABLE IF NOT EXISTS pull_requests (
  pull_request_id TEXT PRIMARY KEY,
  pull_request_name TEXT NOT NULL,
  author_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
  status TEXT NOT NULL DEFAULT 'OPEN',
  need_more_reviewers BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  merged_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_pr_author ON pull_requests(author_id);
CREATE INDEX IF NOT EXISTS idx_pr_status ON pull_requests(status);

-- pr reviewers (assignments) — up to 2 per PR, but DB doesn't enforce numerically
CREATE TABLE IF NOT EXISTS pr_reviewers (
  pr_id TEXT NOT NULL REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
  user_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
  assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (pr_id, user_id)
);

COMMIT;
