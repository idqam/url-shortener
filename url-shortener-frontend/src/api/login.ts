import { supabase } from "../lib/supabaseClient";

export async function loginUser(email: string, password: string) {
  const { data, error } = await supabase.auth.signInWithPassword({
    email,
    password,
  });

  if (error || !data.session) {
    throw new Error(error?.message || "Login failed");
  }

  return data.session.user.id;
}
