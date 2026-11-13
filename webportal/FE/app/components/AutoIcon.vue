<script setup>
import { computed } from 'vue'

const props = defineProps({
  text: { type: String, default: '' },
  size: { type: String, default: 'w-5 h-5' }
})

const iconName = computed(() => {
  const t = props.text.toLowerCase()
  if (t.includes('didik') || t.includes('sekolah') || t.includes('siswa') || t.includes('beasiswa')) return 'i-tabler-school'
  if (t.includes('sehat') || t.includes('rumah sakit') || t.includes('puskesmas') || t.includes('medis')) return 'i-tabler-heartbeat'
  if (t.includes('pajak') || t.includes('retribusi') || t.includes('pbb') || t.includes('samsat')) return 'i-tabler-receipt'
  if (t.includes('sipil') || t.includes('penduduk') || t.includes('ktp') || t.includes('kk') || t.includes('akta')) return 'i-tabler-users'
  if (t.includes('izin') || t.includes('usaha') || t.includes('niaga') || t.includes('perizinan')) return 'i-tabler-file-certificate'
  if (t.includes('sosial') || t.includes('bantuan') || t.includes('bansos')) return 'i-tabler-heart-handshake'
  if (t.includes('lingkungan') || t.includes('sampah') || t.includes('taman')) return 'i-tabler-leaf'
  if (t.includes('aduan') || t.includes('lapor') || t.includes('keluhan')) return 'i-tabler-speakerphone'
  if (t.includes('kerja') || t.includes('tenaga') || t.includes('loker')) return 'i-tabler-briefcase'
  if (t.includes('bangunan') || t.includes('infrastruktur') || t.includes('jalan')) return 'i-tabler-building'
  if (t.includes('data') || t.includes('statistik') || t.includes('informasi')) return 'i-tabler-chart-pie'
  if (t.includes('hukum') || t.includes('perda')) return 'i-tabler-gavel'
  if (t.includes('ambulans') || t.includes('tems') || t.includes('medis')) return 'i-tabler-ambulance'
  if (t.includes('kebakaran') || t.includes('damkar') || t.includes('api')) return 'i-tabler-flame'
  if (t.includes('polisi') || t.includes('kepolisian')) return 'i-tabler-shield'
  if (t.includes('sar') || t.includes('penyelamatan')) return 'i-tabler-lifebuoy'
  if (t.includes('pln') || t.includes('listrik')) return 'i-tabler-bolt'
  if (t.includes('pdam') || t.includes('air')) return 'i-tabler-droplet'
  if (t.includes('bpbd') || t.includes('bencana')) return 'i-tabler-alert-triangle'
  if (t.includes('wisata') || t.includes('pariwisata') || t.includes('destinasi') || t.includes('rekreasi')) return 'i-tabler-map-pin-filled'
  if (t.includes('tani') || t.includes('pertanian') || t.includes('pangan') || t.includes('panen')) return 'i-tabler-plant'
  if (t.includes('ikan') || t.includes('perikanan') || t.includes('nelayan') || t.includes('bahari')) return 'i-tabler-fish'
  if (t.includes('ternak') || t.includes('peternakan') || t.includes('hewan')) return 'i-mdi-cow'
  if (t.includes('umkm') || t.includes('pasar') || t.includes('dagang') || t.includes('perdagangan') || t.includes('toko')) return 'i-tabler-shopping-bag'
  if (t.includes('dishub') || t.includes('transportasi') || t.includes('parkir') || t.includes('lalu lintas') || t.includes('trayek')) return 'i-tabler-bus'
  if (t.includes('budaya') || t.includes('seni') || t.includes('kesenian') || t.includes('museum') || t.includes('cagar')) return 'i-mdi-palette'
  if (t.includes('pemuda') || t.includes('olahraga') || t.includes('dispora') || t.includes('atlet')) return 'i-tabler-trophy'
  if (t.includes('berita') || t.includes('pengumuman') || t.includes('warta') || t.includes('kabar')) return 'i-tabler-news'
  if (t.includes('agenda') || t.includes('acara') || t.includes('kegiatan') || t.includes('jadwal')) return 'i-tabler-calendar-event'
  if (t.includes('galeri') || t.includes('foto') || t.includes('video') || t.includes('dokumentasi')) return 'i-tabler-photo'
  if (t.includes('apbd') || t.includes('anggaran') || t.includes('keuangan') || t.includes('kas')) return 'i-tabler-cash'
  if (t.includes('transparansi') || t.includes('sakip') || t.includes('lakip') || t.includes('laporan')) return 'i-tabler-report-analytics'
  if (t.includes('lpse') || t.includes('tender') || t.includes('pengadaan') || t.includes('pbj')) return 'i-mdi-cart-check'
  if (t.includes('opd') || t.includes('dinas') || t.includes('organisasi') || t.includes('pejabat')) return 'i-tabler-building-community'
  if (t.includes('posyandu') || t.includes('stunting') || t.includes('kb') || t.includes('balita')) return 'i-tabler-baby-carriage'
  if (t.includes('download') || t.includes('unduh') || t.includes('dokumen') || t.includes('file')) return 'i-tabler-file-download'
  if (t.includes('peta') || t.includes('gis') || t.includes('geospasial') || t.includes('pemetaan')) return 'i-tabler-map-2'
  if (t.includes('kontak') || t.includes('hotline') || t.includes('call center') || t.includes('telepon')) return 'i-tabler-phone-call'
  if (t.includes('kerjasama') || t.includes('mou') || t.includes('kemitraan')) return 'i-tabler-handshake'
  if (t.includes('satpol') || t.includes('ketertiban') || t.includes('trantib') || t.includes('penertiban')) return 'i-mdi-police-badge'
  if (t.includes('perpus') || t.includes('perpustakaan') || t.includes('literasi') || t.includes('arsip')) return 'i-tabler-books'
  if (t.includes('inovasi') || t.includes('riset') || t.includes('litbang') || t.includes('penelitian')) return 'i-tabler-bulb'
  if (t.includes('kesra') || t.includes('agama') || t.includes('keagamaan') || t.includes('zakat') || t.includes('hibah')) return 'i-mdi-mosque'
  if (t.includes('perempuan') || t.includes('anak') || t.includes('gender') || t.includes('pppa')) return 'i-tabler-woman'
  if (t.includes('spbe') || t.includes('siber') || t.includes('sandi') || t.includes('aplikasi') || t.includes('website')) return 'i-tabler-devices-pc'
  if (t.includes('bpjs') || t.includes('jkn') || t.includes('kis') || t.includes('jamsostek')) return 'i-mdi-card-account-details'
  if (t.includes('pelatihan') || t.includes('blk') || t.includes('kursus') || t.includes('keterampilan')) return 'i-tabler-certificate-2'
  if (t.includes('pertanahan') || t.includes('tata ruang') || t.includes('sertifikat') || t.includes('tanah')) return 'i-tabler-ruler-measure'
  if (t.includes('surat') || t.includes('disposisi') || t.includes('e-office') || t.includes('persuratan')) return 'i-tabler-mail'
  if (t.includes('daftar') || t.includes('pendaftaran') || t.includes('formulir') || t.includes('registrasi')) return 'i-tabler-forms'
  if (t.includes('bumd') || t.includes('perumda') || t.includes('bpr')) return 'i-tabler-building-bank'
  if (t.includes('halal') || t.includes('sertifikasi') || t.includes('mutu') || t.includes('iso')) return 'i-mdi-certificate'
  if (t.includes('faq') || t.includes('sop') || t.includes('panduan') || t.includes('petunjuk') || t.includes('syarat')) return 'i-tabler-help-circle'
  if (t.includes('inflasi') || t.includes('tpid') || t.includes('harga') || t.includes('sembako')) return 'i-tabler-scale-outline'
  if (t.includes('disabilitas') || t.includes('infrastruktur inklusif') || t.includes('difabel')) return 'i-tabler-disabled'
  return 'i-tabler-file-text' // Default
})
</script>

<template>
  <UIcon :name="iconName" :class="size" />
</template>