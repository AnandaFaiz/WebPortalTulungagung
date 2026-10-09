<script setup>
const links = [
    { label: 'Dashboard', icon: 'i-tabler-layout-dashboard', to: '/admin' },
    { label: 'CMS', icon: 'i-tabler-file-text', to: '/admin/cms' },
    { label: 'Layanan Publik', icon: 'i-tabler-briefcase', to: '/admin/layanan' },
    { label: 'Pengaduan', icon: 'i-tabler-message', to: '/admin/pengaduan' },
]
</script>

<template>
    <UDashboardGroup>
        <UDashboardSidebar collapsible resizable>
            <template #header="{ collapsed }">
                <div class="flex items-center gap-2">
                    <UIcon name="i-tabler-building-bank" class="text-primary text-2xl" />
                    <span v-if="!collapsed" class="font-bold text-lg">Admin Panel</span>
                </div>
            </template>

            <UNavigationMenu :items="links" orientation="vertical" />
        </UDashboardSidebar>

        <UDashboardPanel>
            <UDashboardNavbar title="Dashboard Admin">
                <template #right>
                    <UButton icon="i-tabler-bell" color="neutral" variant="ghost" />
                    <UButton icon="i-tabler-user" color="neutral" variant="ghost" />
                </template>
            </UDashboardNavbar>

            <UDashboardPanelContent>
                <slot />
            </UDashboardPanelContent>
        </UDashboardPanel>
    </UDashboardGroup>
</template>