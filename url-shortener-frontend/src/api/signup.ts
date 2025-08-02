import { supabase } from "../lib/supabaseClient";

export async function signupAndRegister(email: string, password: string) {
  const { error: signupError } = await supabase.auth.signUp({
    email,
    password,
  });
  if (signupError) throw new Error(signupError.message);

  const { data, error: loginError } = await supabase.auth.signInWithPassword({
    email,
    password,
  });
  if (loginError || !data.session)
    throw new Error(loginError?.message || "Login failed");

  return data.session.user.id;
}
