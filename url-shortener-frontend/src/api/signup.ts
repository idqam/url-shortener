import { supabase } from "../lib/supabaseClient";

export async function signupAndRegister(email: string, password: string) {
  const { data: signupData, error: signupError } = await supabase.auth.signUp({
    email,
    password,
  });

  if (signupError) throw new Error(signupError.message);

  const user = signupData.user;
  const session = signupData.session;

  if (!session || !user) {
    console.log("SignUp result:", signupError);
    throw new Error(
      "Sign-up failed to return a session. Try logging in manually."
    );
  }

  return user.id;
}
