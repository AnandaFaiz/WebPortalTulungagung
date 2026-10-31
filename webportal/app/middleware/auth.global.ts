import { getRequiredPermission } from "~/utils/permissions";

export default defineNuxtRouteMiddleware(async (to) => {
  const { user, initialized, fetchUser, can } = useAuth();

  // Tunggu sampai user selesai di-fetch

  if (!initialized.value) {
    await fetchUser();
  }

  // ==== Proteksi /admin/* ====
  if (to.path.startsWith("/admin")) {
    if (!user.value) {
      return navigateTo("/login");
    }
    const requiredPermission = getRequiredPermission(to.path);
    if (requiredPermission && !can(requiredPermission)) {
      return navigateTo("/admin?unauthorized=1");
    }
  }
});
