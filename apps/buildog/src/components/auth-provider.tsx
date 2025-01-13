import React, { createContext, useContext, useEffect, useState } from "react";
import { usePathname, useRouter } from "next/navigation";
import { Auth, User, Service } from "@/web-sdk";

interface UserInformation {
  first_name: string;
  last_name: string;
  email: string;
}

interface AuthStateType {
  user: User | null;
  userInformation: UserInformation | null;
  loading: boolean;
}

const AuthContext = createContext<AuthStateType>({
  user: null,
  userInformation: null,
  loading: true,
});

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [authState, setAuthState] = useState<AuthStateType>({
    user: null,
    userInformation: null,
    loading: true,
  });
  const router = useRouter();
  const pathname = usePathname();

  const fetchUserInformation = async (user: User) => {
    try {
      const response = await Service.makeAuthenticatedRequest("user");
      if (!response.error) {
        setAuthState((prev) => ({
          ...prev,
          userInformation: response,
          loading: false,
        }));
      }
    } catch (error) {
      console.error("Failed to fetch user credentials:", error);
      setAuthState((prev) => ({ ...prev, loading: false }));
    }
  };

  useEffect(() => {
    const handleAuthStateChange = async (user: User | null) => {
      if (user) {
        if (!user.emailVerified && (pathname === "/login/" || pathname === "/signup/")) {
          return;
        }

        setAuthState((prev) => ({ ...prev, user, loading: true }));
        await fetchUserInformation(user);

        // Redirect based on sign-in state and current pathname
        if (pathname === "/login" || pathname === "/signup") {
          router.push("/organizations/");
        }
        return;
      }

      setAuthState({ user: null, userInformation: null, loading: false });
      if (!(pathname === "/login" || pathname === "/signup")) {
        router.push("/login");
      }
    };

    // Listen for auth state changes
    const listen = Auth.onAuthStateChange(handleAuthStateChange);

    // Clean up the auth state listener on unmount
    return () => listen();
  }, [router, pathname]);

  return <AuthContext.Provider value={authState}>{children}</AuthContext.Provider>;
};

export const useAuth = () => useContext(AuthContext);
