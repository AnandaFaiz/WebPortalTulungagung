export interface User {
  id: number;
  nama: string;
  email: string;
  role: "super_admin" | "admin_opd" | "editor";
  opd: string | null;
  is_active?: boolean;
  created_at?: string;
  last_login_at?: string | null;
}

const API_BASE = "http://localhost:8080/api";

export function useAuth() {
  const user = useState<User | null>("auth_user", () => null);
  const permissions = useState<Record<string, boolean>>(
    "auth_permissions",
    () => ({}),
  );
  const loading = useState<boolean>("auth_loading", () => false);
  const initialized = useState<boolean>("auth_initialized", () => false);

  const login = async (email: string, password: string) => {
    loading.value = true;
    try {
      const response = await $fetch<{ success: boolean; user: User }>(
        `${API_BASE}/auth/login`,
        {
          method: "POST",
          body: { email, password },
          credentials: "include",
        },
      );
      user.value = response.user;
      await fetchPermissions();
      initialized.value = true;

      if (import.meta.client) {
        sessionStorage.setItem("auth_tab_active", "true");
      }

      return { success: true, user: response.user };
    } catch (err: any) {
      return {
        success: false,
        message: err?.data?.statusMessage || "Login gagal.",
      };
    } finally {
      loading.value = false;
    }
  };

  const fetchUser = async (): Promise<User | null> => {
    try {
      const requestFetch = import.meta.server ? useRequestFetch() : $fetch;
      const response = await requestFetch<{ user: User }>(
        `${API_BASE}/auth/me`,
        { credentials: "include" },
      );
      user.value = response.user;
      await fetchPermissions();
      initialized.value = true;
      return response.user;
    } catch (err) {
      user.value = null;
      permissions.value = {};
      initialized.value = true;
      return null;
    }
  };

  const logout = async () => {
    try {
      await $fetch(`${API_BASE}/auth/logout`, {
        method: "POST",
        credentials: "include",
      });
    } catch {}
    user.value = null;
    permissions.value = {};
    initialized.value = false;
    if (import.meta.client) {
      sessionStorage.removeItem("auth_tab_active");
    }
    await navigateTo("/login");
  };
  const fetchPermissions = async () => {
    try {
      const requestFetch = import.meta.server ? useRequestFetch() : $fetch;

      const response = await requestFetch<{
        role: string;
        permissions: Record<string, boolean>;
      }>(`${API_BASE}/auth/permissions`, { credentials: "include" });
      permissions.value = response.permissions;
    } catch {
      permissions.value = {};
    }
  };

  const can = (permissionKey: string): boolean => {
    if (user.value?.role === "super_admin") return true;
    return permissions.value[permissionKey] === true;
  };

  const isLoggedIn = computed(() => !!user.value);
  const isSuperAdmin = computed(() => user.value?.role === "super_admin");
  const isAdminOpd = computed(() => user.value?.role === "admin_opd");
  const isEditor = computed(() => user.value?.role === "editor");

  return {
    user,
    permissions,
    loading,
    initialized,
    login,
    logout,
    fetchUser,
    fetchPermissions,
    can,
    isLoggedIn,
    isSuperAdmin,
    isAdminOpd,
    isEditor,
  };
}
