--
-- PostgreSQL database dump
--

-- Dumped from database version 18.6
-- Dumped by pg_dump version 18.6

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: activity_logs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.activity_logs (
    id integer NOT NULL,
    user_id integer,
    user_nama character varying(100),
    user_role character varying(50),
    action character varying(50) NOT NULL,
    entity character varying(50),
    entity_id integer,
    description text,
    ip_address character varying(45),
    user_agent text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: activity_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.activity_logs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: activity_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.activity_logs_id_seq OWNED BY public.activity_logs.id;


--
-- Name: agenda; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.agenda (
    id integer NOT NULL,
    judul character varying(255) NOT NULL,
    deskripsi text,
    tanggal timestamp without time zone NOT NULL,
    lokasi character varying(255),
    gambar character varying(255),
    kategori character varying(100),
    status character varying(20) DEFAULT 'Publish'::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    link_event character varying(255),
    CONSTRAINT agenda_status_check CHECK (((status)::text = ANY ((ARRAY['Draft'::character varying, 'Publish'::character varying])::text[])))
);


--
-- Name: agenda_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.agenda_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: agenda_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.agenda_id_seq OWNED BY public.agenda.id;


--
-- Name: aplikasi; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.aplikasi (
    id integer NOT NULL,
    nama character varying(150) NOT NULL,
    deskripsi text,
    icon character varying(100) DEFAULT 'i-tabler-apps'::character varying,
    url character varying(255),
    urutan integer DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: aplikasi_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.aplikasi_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: aplikasi_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.aplikasi_id_seq OWNED BY public.aplikasi.id;


--
-- Name: banner; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.banner (
    id integer NOT NULL,
    judul character varying(255) NOT NULL,
    gambar character varying(255),
    link character varying(255),
    urutan integer DEFAULT 0,
    status character varying(20) DEFAULT 'Aktif'::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    tipe character varying(20) DEFAULT 'hero'::character varying,
    CONSTRAINT banner_status_check CHECK (((status)::text = ANY ((ARRAY['Aktif'::character varying, 'Nonaktif'::character varying])::text[])))
);


--
-- Name: banner_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.banner_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: banner_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.banner_id_seq OWNED BY public.banner.id;


--
-- Name: berita; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.berita (
    id integer NOT NULL,
    judul character varying(255) NOT NULL,
    deskripsi text,
    konten text,
    kategori character varying(50) DEFAULT 'Pemerintahan'::character varying,
    penulis character varying(100),
    tanggal date DEFAULT CURRENT_DATE,
    status character varying(20) DEFAULT 'Draft'::character varying,
    thumbnail character varying(255),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    views integer DEFAULT 0,
    CONSTRAINT berita_status_check CHECK (((status)::text = ANY ((ARRAY['Draft'::character varying, 'Publish'::character varying])::text[])))
);


--
-- Name: berita_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.berita_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: berita_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.berita_id_seq OWNED BY public.berita.id;


--
-- Name: galeri; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.galeri (
    id integer NOT NULL,
    judul character varying(255) NOT NULL,
    deskripsi text,
    tipe character varying(20) DEFAULT 'foto'::character varying,
    tanggal date DEFAULT CURRENT_DATE,
    status character varying(20) DEFAULT 'Draft'::character varying,
    thumbnail character varying(255),
    file_url character varying(255),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT galeri_status_check CHECK (((status)::text = ANY ((ARRAY['Draft'::character varying, 'Publish'::character varying])::text[]))),
    CONSTRAINT galeri_tipe_check CHECK (((tipe)::text = ANY ((ARRAY['foto'::character varying, 'video'::character varying])::text[])))
);


--
-- Name: galeri_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.galeri_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: galeri_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.galeri_id_seq OWNED BY public.galeri.id;


--
-- Name: halaman_statis; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.halaman_statis (
    id integer NOT NULL,
    judul character varying(255) NOT NULL,
    slug character varying(255) NOT NULL,
    konten text,
    tanggal date DEFAULT CURRENT_DATE,
    status character varying(20) DEFAULT 'Draft'::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT halaman_statis_status_check CHECK (((status)::text = ANY ((ARRAY['Draft'::character varying, 'Publish'::character varying])::text[])))
);


--
-- Name: halaman_statis_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.halaman_statis_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: halaman_statis_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.halaman_statis_id_seq OWNED BY public.halaman_statis.id;


--
-- Name: identitas_situs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.identitas_situs (
    id integer NOT NULL,
    nama_situs character varying(150) DEFAULT 'Portal Pemkab Tulungagung'::character varying,
    tagline character varying(255) DEFAULT 'Portal Resmi Pemerintah Daerah'::character varying,
    deskripsi text,
    logo character varying(255) DEFAULT '/Logo-Kabupaten-Tulungagung.png'::character varying,
    favicon character varying(255),
    email character varying(150),
    telepon character varying(50),
    alamat text,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: identitas_situs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.identitas_situs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: identitas_situs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.identitas_situs_id_seq OWNED BY public.identitas_situs.id;


--
-- Name: informasi_mpp; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.informasi_mpp (
    id integer NOT NULL,
    deskripsi text,
    alamat character varying(255),
    jam_operasional character varying(100),
    jumlah_layanan integer DEFAULT 0,
    jumlah_instansi integer DEFAULT 0,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: informasi_mpp_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.informasi_mpp_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: informasi_mpp_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.informasi_mpp_id_seq OWNED BY public.informasi_mpp.id;


--
-- Name: kategori_berita; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.kategori_berita (
    id integer NOT NULL,
    nama character varying(100) NOT NULL,
    icon character varying(100) DEFAULT 'i-tabler-news'::character varying,
    deskripsi text
);


--
-- Name: kategori_berita_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.kategori_berita_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: kategori_berita_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.kategori_berita_id_seq OWNED BY public.kategori_berita.id;


--
-- Name: kategori_layanan; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.kategori_layanan (
    id integer NOT NULL,
    nama character varying(100) NOT NULL,
    deskripsi text,
    icon character varying(100) DEFAULT 'i-tabler-file-text'::character varying,
    urutan integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: kategori_layanan_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.kategori_layanan_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: kategori_layanan_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.kategori_layanan_id_seq OWNED BY public.kategori_layanan.id;


--
-- Name: kontak; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.kontak (
    id integer NOT NULL,
    alamat text,
    telepon character varying(50),
    email character varying(150),
    jam_operasional character varying(150),
    embed_maps text,
    facebook character varying(255),
    instagram character varying(255),
    twitter character varying(255),
    youtube character varying(255),
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: kontak_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.kontak_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: kontak_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.kontak_id_seq OWNED BY public.kontak.id;


--
-- Name: layanan; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.layanan (
    id integer NOT NULL,
    kategori_id integer,
    nama character varying(150) NOT NULL,
    link character varying(255),
    urutan integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: layanan_darurat; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.layanan_darurat (
    id integer NOT NULL,
    nama character varying(100) NOT NULL,
    nomor character varying(50) NOT NULL,
    icon character varying(100) DEFAULT 'i-tabler-phone'::character varying,
    urutan integer DEFAULT 0,
    is_active boolean DEFAULT true
);


--
-- Name: layanan_darurat_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.layanan_darurat_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: layanan_darurat_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.layanan_darurat_id_seq OWNED BY public.layanan_darurat.id;


--
-- Name: layanan_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.layanan_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: layanan_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.layanan_id_seq OWNED BY public.layanan.id;


--
-- Name: layanan_terpopuler; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.layanan_terpopuler (
    id integer NOT NULL,
    kategori character varying(50) NOT NULL,
    judul character varying(150) NOT NULL,
    deskripsi text,
    icon character varying(100),
    link character varying(255),
    urutan integer DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: layanan_terpopuler_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.layanan_terpopuler_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: layanan_terpopuler_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.layanan_terpopuler_id_seq OWNED BY public.layanan_terpopuler.id;


--
-- Name: media; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.media (
    id integer NOT NULL,
    nama_file character varying(255) NOT NULL,
    url character varying(255) NOT NULL,
    tipe_file character varying(50),
    ukuran_file integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: media_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.media_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: media_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.media_id_seq OWNED BY public.media.id;


--
-- Name: menu_navigasi; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.menu_navigasi (
    id integer NOT NULL,
    label character varying(120) NOT NULL,
    url character varying(255) DEFAULT '#'::character varying NOT NULL,
    deskripsi text DEFAULT ''::text,
    icon character varying(120) DEFAULT 'i-tabler-link'::character varying,
    color character varying(80) DEFAULT 'text-sky-600'::character varying,
    grup character varying(120) DEFAULT 'Lainnya'::character varying NOT NULL,
    urutan integer DEFAULT 0 NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


--
-- Name: menu_navigasi_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.menu_navigasi_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: menu_navigasi_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.menu_navigasi_id_seq OWNED BY public.menu_navigasi.id;


--
-- Name: opds; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.opds (
    id integer NOT NULL,
    nama character varying(150) NOT NULL,
    deskripsi text,
    website character varying(255),
    urutan integer DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: opds_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.opds_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: opds_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.opds_id_seq OWNED BY public.opds.id;


--
-- Name: pejabat; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.pejabat (
    id integer NOT NULL,
    nama character varying(150) NOT NULL,
    jabatan character varying(150) NOT NULL,
    kategori character varying(100),
    foto character varying(255),
    periode character varying(50),
    urutan integer DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: pejabat_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.pejabat_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: pejabat_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.pejabat_id_seq OWNED BY public.pejabat.id;


--
-- Name: pencarian_populer; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.pencarian_populer (
    id integer NOT NULL,
    keyword character varying(100) NOT NULL,
    urutan integer DEFAULT 0
);


--
-- Name: pencarian_populer_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.pencarian_populer_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: pencarian_populer_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.pencarian_populer_id_seq OWNED BY public.pencarian_populer.id;


--
-- Name: pengaduan; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.pengaduan (
    id integer NOT NULL,
    nama character varying(100) NOT NULL,
    kontak character varying(100) NOT NULL,
    kategori character varying(100) NOT NULL,
    isi_aduan text NOT NULL,
    status character varying(20) DEFAULT 'Baru'::character varying,
    balasan_admin text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pengaduan_status_check CHECK (((status)::text = ANY ((ARRAY['Baru'::character varying, 'Diproses'::character varying, 'Selesai'::character varying])::text[])))
);


--
-- Name: pengaduan_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.pengaduan_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: pengaduan_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.pengaduan_id_seq OWNED BY public.pengaduan.id;


--
-- Name: pengaduan_kategori; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.pengaduan_kategori (
    id integer NOT NULL,
    nama character varying(120) NOT NULL,
    deskripsi text DEFAULT ''::text,
    urutan integer DEFAULT 0 NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


--
-- Name: pengaduan_kategori_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.pengaduan_kategori_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: pengaduan_kategori_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.pengaduan_kategori_id_seq OWNED BY public.pengaduan_kategori.id;


--
-- Name: pengaduan_status; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.pengaduan_status (
    id integer NOT NULL,
    nama character varying(120) NOT NULL,
    deskripsi text DEFAULT ''::text,
    urutan integer DEFAULT 0 NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);


--
-- Name: pengaduan_status_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.pengaduan_status_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: pengaduan_status_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.pengaduan_status_id_seq OWNED BY public.pengaduan_status.id;


--
-- Name: pengumuman; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.pengumuman (
    id integer NOT NULL,
    judul character varying(255) NOT NULL,
    konten text,
    tanggal date DEFAULT CURRENT_DATE,
    prioritas character varying(20) DEFAULT 'Normal'::character varying,
    status character varying(20) DEFAULT 'Draft'::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pengumuman_prioritas_check CHECK (((prioritas)::text = ANY ((ARRAY['Normal'::character varying, 'Tinggi'::character varying])::text[]))),
    CONSTRAINT pengumuman_status_check CHECK (((status)::text = ANY ((ARRAY['Draft'::character varying, 'Publish'::character varying])::text[])))
);


--
-- Name: pengumuman_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.pengumuman_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: pengumuman_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.pengumuman_id_seq OWNED BY public.pengumuman.id;


--
-- Name: profil_daerah; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.profil_daerah (
    id integer NOT NULL,
    sejarah text,
    lambang text,
    visi text,
    misi jsonb DEFAULT '[]'::jsonb,
    geografis text,
    demografi jsonb DEFAULT '{}'::jsonb,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: profil_daerah_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.profil_daerah_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: profil_daerah_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.profil_daerah_id_seq OWNED BY public.profil_daerah.id;


--
-- Name: role_permissions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.role_permissions (
    id integer NOT NULL,
    role character varying(50) NOT NULL,
    permission_key character varying(100) NOT NULL,
    allowed boolean DEFAULT false,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: role_permissions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.role_permissions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: role_permissions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.role_permissions_id_seq OWNED BY public.role_permissions.id;


--
-- Name: struktur_organisasi; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.struktur_organisasi (
    id integer NOT NULL,
    judul character varying(255) DEFAULT 'Struktur Organisasi Pemkab Tulungagung'::character varying,
    deskripsi text,
    gambar character varying(255),
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: struktur_organisasi_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.struktur_organisasi_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: struktur_organisasi_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.struktur_organisasi_id_seq OWNED BY public.struktur_organisasi.id;


--
-- Name: tautan_eksternal; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tautan_eksternal (
    id integer NOT NULL,
    nama character varying(150) NOT NULL,
    url character varying(255) NOT NULL,
    deskripsi text,
    icon character varying(100) DEFAULT 'i-tabler-external-link'::character varying,
    urutan integer DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: tautan_eksternal_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tautan_eksternal_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tautan_eksternal_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tautan_eksternal_id_seq OWNED BY public.tautan_eksternal.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id integer NOT NULL,
    nama character varying(100) NOT NULL,
    email character varying(150) NOT NULL,
    password character varying(255) NOT NULL,
    role character varying(50) NOT NULL,
    opd character varying(100),
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_role_check CHECK (((role)::text = ANY ((ARRAY['super_admin'::character varying, 'admin_opd'::character varying, 'editor'::character varying])::text[])))
);


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: activity_logs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.activity_logs ALTER COLUMN id SET DEFAULT nextval('public.activity_logs_id_seq'::regclass);


--
-- Name: agenda id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.agenda ALTER COLUMN id SET DEFAULT nextval('public.agenda_id_seq'::regclass);


--
-- Name: aplikasi id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.aplikasi ALTER COLUMN id SET DEFAULT nextval('public.aplikasi_id_seq'::regclass);


--
-- Name: banner id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.banner ALTER COLUMN id SET DEFAULT nextval('public.banner_id_seq'::regclass);


--
-- Name: berita id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.berita ALTER COLUMN id SET DEFAULT nextval('public.berita_id_seq'::regclass);


--
-- Name: galeri id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.galeri ALTER COLUMN id SET DEFAULT nextval('public.galeri_id_seq'::regclass);


--
-- Name: halaman_statis id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.halaman_statis ALTER COLUMN id SET DEFAULT nextval('public.halaman_statis_id_seq'::regclass);


--
-- Name: identitas_situs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.identitas_situs ALTER COLUMN id SET DEFAULT nextval('public.identitas_situs_id_seq'::regclass);


--
-- Name: informasi_mpp id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.informasi_mpp ALTER COLUMN id SET DEFAULT nextval('public.informasi_mpp_id_seq'::regclass);


--
-- Name: kategori_berita id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.kategori_berita ALTER COLUMN id SET DEFAULT nextval('public.kategori_berita_id_seq'::regclass);


--
-- Name: kategori_layanan id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.kategori_layanan ALTER COLUMN id SET DEFAULT nextval('public.kategori_layanan_id_seq'::regclass);


--
-- Name: kontak id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.kontak ALTER COLUMN id SET DEFAULT nextval('public.kontak_id_seq'::regclass);


--
-- Name: layanan id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.layanan ALTER COLUMN id SET DEFAULT nextval('public.layanan_id_seq'::regclass);


--
-- Name: layanan_darurat id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.layanan_darurat ALTER COLUMN id SET DEFAULT nextval('public.layanan_darurat_id_seq'::regclass);


--
-- Name: layanan_terpopuler id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.layanan_terpopuler ALTER COLUMN id SET DEFAULT nextval('public.layanan_terpopuler_id_seq'::regclass);


--
-- Name: media id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.media ALTER COLUMN id SET DEFAULT nextval('public.media_id_seq'::regclass);


--
-- Name: menu_navigasi id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.menu_navigasi ALTER COLUMN id SET DEFAULT nextval('public.menu_navigasi_id_seq'::regclass);


--
-- Name: opds id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.opds ALTER COLUMN id SET DEFAULT nextval('public.opds_id_seq'::regclass);


--
-- Name: pejabat id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pejabat ALTER COLUMN id SET DEFAULT nextval('public.pejabat_id_seq'::regclass);


--
-- Name: pencarian_populer id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pencarian_populer ALTER COLUMN id SET DEFAULT nextval('public.pencarian_populer_id_seq'::regclass);


--
-- Name: pengaduan id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pengaduan ALTER COLUMN id SET DEFAULT nextval('public.pengaduan_id_seq'::regclass);


--
-- Name: pengaduan_kategori id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pengaduan_kategori ALTER COLUMN id SET DEFAULT nextval('public.pengaduan_kategori_id_seq'::regclass);


--
-- Name: pengaduan_status id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pengaduan_status ALTER COLUMN id SET DEFAULT nextval('public.pengaduan_status_id_seq'::regclass);


--
-- Name: pengumuman id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pengumuman ALTER COLUMN id SET DEFAULT nextval('public.pengumuman_id_seq'::regclass);


--
-- Name: profil_daerah id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.profil_daerah ALTER COLUMN id SET DEFAULT nextval('public.profil_daerah_id_seq'::regclass);


--
-- Name: role_permissions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_permissions ALTER COLUMN id SET DEFAULT nextval('public.role_permissions_id_seq'::regclass);


--
-- Name: struktur_organisasi id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.struktur_organisasi ALTER COLUMN id SET DEFAULT nextval('public.struktur_organisasi_id_seq'::regclass);


--
-- Name: tautan_eksternal id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tautan_eksternal ALTER COLUMN id SET DEFAULT nextval('public.tautan_eksternal_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: activity_logs activity_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.activity_logs
    ADD CONSTRAINT activity_logs_pkey PRIMARY KEY (id);


--
-- Name: agenda agenda_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.agenda
    ADD CONSTRAINT agenda_pkey PRIMARY KEY (id);


--
-- Name: aplikasi aplikasi_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.aplikasi
    ADD CONSTRAINT aplikasi_pkey PRIMARY KEY (id);


--
-- Name: banner banner_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.banner
    ADD CONSTRAINT banner_pkey PRIMARY KEY (id);


--
-- Name: berita berita_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.berita
    ADD CONSTRAINT berita_pkey PRIMARY KEY (id);


--
-- Name: galeri galeri_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.galeri
    ADD CONSTRAINT galeri_pkey PRIMARY KEY (id);


--
-- Name: halaman_statis halaman_statis_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.halaman_statis
    ADD CONSTRAINT halaman_statis_pkey PRIMARY KEY (id);


--
-- Name: halaman_statis halaman_statis_slug_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.halaman_statis
    ADD CONSTRAINT halaman_statis_slug_key UNIQUE (slug);


--
-- Name: identitas_situs identitas_situs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.identitas_situs
    ADD CONSTRAINT identitas_situs_pkey PRIMARY KEY (id);


--
-- Name: informasi_mpp informasi_mpp_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.informasi_mpp
    ADD CONSTRAINT informasi_mpp_pkey PRIMARY KEY (id);


--
-- Name: kategori_berita kategori_berita_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.kategori_berita
    ADD CONSTRAINT kategori_berita_pkey PRIMARY KEY (id);


--
-- Name: kategori_layanan kategori_layanan_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.kategori_layanan
    ADD CONSTRAINT kategori_layanan_pkey PRIMARY KEY (id);


--
-- Name: kontak kontak_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.kontak
    ADD CONSTRAINT kontak_pkey PRIMARY KEY (id);


--
-- Name: layanan_darurat layanan_darurat_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.layanan_darurat
    ADD CONSTRAINT layanan_darurat_pkey PRIMARY KEY (id);


--
-- Name: layanan layanan_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.layanan
    ADD CONSTRAINT layanan_pkey PRIMARY KEY (id);


--
-- Name: layanan_terpopuler layanan_terpopuler_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.layanan_terpopuler
    ADD CONSTRAINT layanan_terpopuler_pkey PRIMARY KEY (id);


--
-- Name: media media_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.media
    ADD CONSTRAINT media_pkey PRIMARY KEY (id);


--
-- Name: menu_navigasi menu_navigasi_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.menu_navigasi
    ADD CONSTRAINT menu_navigasi_pkey PRIMARY KEY (id);


--
-- Name: opds opds_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.opds
    ADD CONSTRAINT opds_pkey PRIMARY KEY (id);


--
-- Name: pejabat pejabat_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pejabat
    ADD CONSTRAINT pejabat_pkey PRIMARY KEY (id);


--
-- Name: pencarian_populer pencarian_populer_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pencarian_populer
    ADD CONSTRAINT pencarian_populer_pkey PRIMARY KEY (id);


--
-- Name: pengaduan_kategori pengaduan_kategori_nama_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pengaduan_kategori
    ADD CONSTRAINT pengaduan_kategori_nama_key UNIQUE (nama);


--
-- Name: pengaduan_kategori pengaduan_kategori_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pengaduan_kategori
    ADD CONSTRAINT pengaduan_kategori_pkey PRIMARY KEY (id);


--
-- Name: pengaduan pengaduan_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pengaduan
    ADD CONSTRAINT pengaduan_pkey PRIMARY KEY (id);


--
-- Name: pengaduan_status pengaduan_status_nama_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pengaduan_status
    ADD CONSTRAINT pengaduan_status_nama_key UNIQUE (nama);


--
-- Name: pengaduan_status pengaduan_status_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pengaduan_status
    ADD CONSTRAINT pengaduan_status_pkey PRIMARY KEY (id);


--
-- Name: pengumuman pengumuman_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pengumuman
    ADD CONSTRAINT pengumuman_pkey PRIMARY KEY (id);


--
-- Name: profil_daerah profil_daerah_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.profil_daerah
    ADD CONSTRAINT profil_daerah_pkey PRIMARY KEY (id);


--
-- Name: role_permissions role_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_pkey PRIMARY KEY (id);


--
-- Name: role_permissions role_permissions_role_permission_key_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_role_permission_key_key UNIQUE (role, permission_key);


--
-- Name: struktur_organisasi struktur_organisasi_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.struktur_organisasi
    ADD CONSTRAINT struktur_organisasi_pkey PRIMARY KEY (id);


--
-- Name: tautan_eksternal tautan_eksternal_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tautan_eksternal
    ADD CONSTRAINT tautan_eksternal_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_activity_logs_action; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_activity_logs_action ON public.activity_logs USING btree (action);


--
-- Name: idx_activity_logs_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_activity_logs_created ON public.activity_logs USING btree (created_at DESC);


--
-- Name: idx_activity_logs_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_activity_logs_user ON public.activity_logs USING btree (user_id);


--
-- Name: idx_aplikasi_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_aplikasi_active ON public.aplikasi USING btree (is_active);


--
-- Name: idx_aplikasi_urutan; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_aplikasi_urutan ON public.aplikasi USING btree (urutan);


--
-- Name: idx_banner_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_banner_status ON public.banner USING btree (status);


--
-- Name: idx_banner_urutan; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_banner_urutan ON public.banner USING btree (urutan);


--
-- Name: idx_berita_kategori; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_berita_kategori ON public.berita USING btree (kategori);


--
-- Name: idx_berita_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_berita_status ON public.berita USING btree (status);


--
-- Name: idx_berita_tanggal; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_berita_tanggal ON public.berita USING btree (tanggal DESC);


--
-- Name: idx_galeri_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_galeri_status ON public.galeri USING btree (status);


--
-- Name: idx_galeri_tanggal; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_galeri_tanggal ON public.galeri USING btree (tanggal DESC);


--
-- Name: idx_halaman_slug; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_halaman_slug ON public.halaman_statis USING btree (slug);


--
-- Name: idx_halaman_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_halaman_status ON public.halaman_statis USING btree (status);


--
-- Name: idx_opds_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_opds_active ON public.opds USING btree (is_active);


--
-- Name: idx_opds_urutan; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_opds_urutan ON public.opds USING btree (urutan);


--
-- Name: idx_pejabat_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_pejabat_active ON public.pejabat USING btree (is_active);


--
-- Name: idx_pejabat_kategori; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_pejabat_kategori ON public.pejabat USING btree (kategori);


--
-- Name: idx_pejabat_urutan; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_pejabat_urutan ON public.pejabat USING btree (urutan);


--
-- Name: idx_pengumuman_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_pengumuman_status ON public.pengumuman USING btree (status);


--
-- Name: idx_pengumuman_tanggal; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_pengumuman_tanggal ON public.pengumuman USING btree (tanggal DESC);


--
-- Name: idx_role_permissions_role; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_role_permissions_role ON public.role_permissions USING btree (role);


--
-- Name: idx_tautan_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_tautan_active ON public.tautan_eksternal USING btree (is_active);


--
-- Name: idx_tautan_urutan; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_tautan_urutan ON public.tautan_eksternal USING btree (urutan);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: activity_logs activity_logs_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.activity_logs
    ADD CONSTRAINT activity_logs_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: layanan layanan_kategori_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.layanan
    ADD CONSTRAINT layanan_kategori_id_fkey FOREIGN KEY (kategori_id) REFERENCES public.kategori_layanan(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--


