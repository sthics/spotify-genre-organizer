-- Fix RLS policies that were too permissive
-- This migration replaces the overly permissive USING (true) policies

-- Drop existing policies
DROP POLICY IF EXISTS "Users can view own data" ON users;
DROP POLICY IF EXISTS "Users can view own jobs" ON organize_jobs;

-- Create proper restrictive policies for users table
CREATE POLICY "Users can view own data" ON users
  FOR SELECT USING (spotify_id = current_setting('app.user_id', true));

CREATE POLICY "Users can update own data" ON users
  FOR UPDATE USING (spotify_id = current_setting('app.user_id', true));

-- Create proper restrictive policies for organize_jobs table
CREATE POLICY "Users can view own jobs" ON organize_jobs
  FOR SELECT USING (user_id = (SELECT id FROM users WHERE spotify_id = current_setting('app.user_id', true)));

CREATE POLICY "Users can insert own jobs" ON organize_jobs
  FOR INSERT WITH CHECK (user_id = (SELECT id FROM users WHERE spotify_id = current_setting('app.user_id', true)));

CREATE POLICY "Users can update own jobs" ON organize_jobs
  FOR UPDATE USING (user_id = (SELECT id FROM users WHERE spotify_id = current_setting('app.user_id', true)));
