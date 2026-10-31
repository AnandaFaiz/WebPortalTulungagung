export default defineNuxtPlugin(async (nuxtApp) => {
  const { user, fetchUser, logout } = useAuth();

  // Cek apakah tab ini punya penanda aktif di sessionStorage
  const isTabActive = sessionStorage.getItem("auth_tab_active");

  // Jika sessionStorage kosong (artinya tab baru / browser baru dibuka)
  if (!isTabActive) {
    // Coba ambil user dari backend (siapa tahu cookie httpOnly masih ada)
    const loggedUser = await fetchUser();

    // Jika backend bilang masih login (cookie ada), TAPI ini tab baru, kita paksa logout!
    if (loggedUser) {
      await logout();
    }
  } else {
    // Jika ini tab yang sama (hanya reload), biarkan saja fetchUser berjalan normal
    if (!user.value) {
      await fetchUser();
    }
  }
});
