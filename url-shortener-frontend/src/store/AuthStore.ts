import { create } from "zustand";
import { supabase } from "../lib/supabaseClient";

type AuthState = {
  userId: string | null;
  accessToken: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  lastUpdate: number; 
  login: (userId: string, accessToken: string) => void;
  logout: () => void;
  setUserIdFromSession: () => Promise<void>;
  getAccessToken: () => string | null;
  setLoading: (loading: boolean) => void;
};

export const useAuthStore = create<AuthState>()((set, get) => ({
  userId: null,
  accessToken: null,
  isAuthenticated: false,
  isLoading: false,
  lastUpdate: 0,

  setLoading: (loading: boolean) => {
    console.log("Setting loading state:", loading);
    set({ isLoading: loading });
  },

  login: (userId, accessToken) => {
    const currentState = get();

    if (
      currentState.userId === userId &&
      currentState.accessToken === accessToken &&
      currentState.isAuthenticated === true
    ) {
      console.log("Login called but state is already correct, skipping update");
      return;
    }

    
    const now = Date.now();
    if (now - currentState.lastUpdate < 100) {
      
      console.log("Login called too soon after last update, debouncing");
      return;
    }

    console.log("Logging in user:", userId);
    set({
      userId,
      accessToken,
      isAuthenticated: true,
      isLoading: false,
      lastUpdate: now,
    });
  },

  logout: () => {
    const currentState = get();

    
    if (
      !currentState.isAuthenticated &&
      !currentState.userId &&
      !currentState.accessToken
    ) {
      console.log("Logout called but user is already logged out, skipping");
      return;
    }

    console.log("Logging out user");

   
    set({
      userId: null,
      accessToken: null,
      isAuthenticated: false,
      isLoading: false,
      lastUpdate: Date.now(),
    });
  },

  setUserIdFromSession: async () => {
    const currentState = get();

    if (currentState.isLoading) {
      console.log("Session check already in progress, skipping");
      return;
    }

    set({ isLoading: true });

    try {
      console.log("Checking session...");
      const { data, error } = await supabase.auth.getSession();

      if (error) {
        console.error("Session check error:", error);
        set({ isLoading: false });
        return;
      }

      const session = data.session;
      if (session?.user?.id && session?.access_token) {
        console.log("Session found, updating state");
        set({
          userId: session.user.id,
          accessToken: session.access_token,
          isAuthenticated: true,
          isLoading: false,
          lastUpdate: Date.now(),
        });
      } else {
        console.log("No session found");
        set({
          userId: null,
          accessToken: null,
          isAuthenticated: false,
          isLoading: false,
          lastUpdate: Date.now(),
        });
      }
    } catch (error) {
      console.error("Unexpected error during session check:", error);
      set({
        isLoading: false,
        lastUpdate: Date.now(),
      });
    }
  },

  getAccessToken: () => {
    const state = get();
    console.log(
      "Getting access token:",
      state.accessToken ? "present" : "missing"
    );
    return state.accessToken;
  },
}));
