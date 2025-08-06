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
    console.log("SignUp result:", signupData);
    throw new Error(
      "Please check your email to confirm your account before logging in."
    );
  }

  return user.id;
}
