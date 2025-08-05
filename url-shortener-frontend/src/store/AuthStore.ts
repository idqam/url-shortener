import { create } from "zustand";
import { supabase } from "../lib/supabaseClient";

type AuthState = {
  userId: string | null;
  accessToken: string | null;
  isAuthenticated: boolean;
  login: (userId: string, accessToken: string) => void;
  logout: () => void;
  setUserIdFromSession: () => Promise<void>;
  getAccessToken: () => string | null;
};

export const useAuthStore = create<AuthState>()((set, get) => ({
  userId: null,
  accessToken: null,
  isAuthenticated: false,

  login: (userId, accessToken) =>
    set({ userId, accessToken, isAuthenticated: true }),

  logout: () => {
    supabase.auth.signOut();
    set({ userId: null, accessToken: null, isAuthenticated: false });
  },

  setUserIdFromSession: async () => {
    const {
      data: { session },
    } = await supabase.auth.getSession();

    if (session?.user?.id && session?.access_token) {
      set({
        userId: session.user.id,
        accessToken: session.access_token,
        isAuthenticated: true,
      });
    } else {
      set({ userId: null, accessToken: null, isAuthenticated: false });
    }
  },

  getAccessToken: () => get().accessToken,
}));
