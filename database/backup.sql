PGDMP     &                
    {         
   backendapi    15.2    15.2 ]    n           0    0    ENCODING    ENCODING        SET client_encoding = 'UTF8';
                      false            o           0    0 
   STDSTRINGS 
   STDSTRINGS     (   SET standard_conforming_strings = 'on';
                      false            p           0    0 
   SEARCHPATH 
   SEARCHPATH     8   SELECT pg_catalog.set_config('search_path', '', false);
                      false            q           1262    16494 
   backendapi    DATABASE     ~   CREATE DATABASE backendapi WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'Russian_Russia.1251';
    DROP DATABASE backendapi;
                postgres    false            �            1259    16519    expense    TABLE     �   CREATE TABLE public.expense (
    id integer NOT NULL,
    amount numeric,
    date date,
    planned boolean,
    user_id integer
);
    DROP TABLE public.expense;
       public         heap    postgres    false            �            1259    16677    expense_categories    TABLE     �   CREATE TABLE public.expense_categories (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    icon character varying(255) NOT NULL,
    is_fixed boolean NOT NULL,
    user_id integer
);
 &   DROP TABLE public.expense_categories;
       public         heap    postgres    false            �            1259    16676    expense_categories_id_seq    SEQUENCE     �   CREATE SEQUENCE public.expense_categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 0   DROP SEQUENCE public.expense_categories_id_seq;
       public          postgres    false    229            r           0    0    expense_categories_id_seq    SEQUENCE OWNED BY     W   ALTER SEQUENCE public.expense_categories_id_seq OWNED BY public.expense_categories.id;
          public          postgres    false    228            �            1259    16518    expense_id_seq    SEQUENCE     �   CREATE SEQUENCE public.expense_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 %   DROP SEQUENCE public.expense_id_seq;
       public          postgres    false    219            s           0    0    expense_id_seq    SEQUENCE OWNED BY     A   ALTER SEQUENCE public.expense_id_seq OWNED BY public.expense.id;
          public          postgres    false    218            �            1259    16559    goal    TABLE     �   CREATE TABLE public.goal (
    id integer NOT NULL,
    goal character varying(255),
    user_id integer,
    need real,
    current_state real
);
    DROP TABLE public.goal;
       public         heap    postgres    false            �            1259    16558    goal_id_seq    SEQUENCE     �   CREATE SEQUENCE public.goal_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 "   DROP SEQUENCE public.goal_id_seq;
       public          postgres    false    223            t           0    0    goal_id_seq    SEQUENCE OWNED BY     ;   ALTER SEQUENCE public.goal_id_seq OWNED BY public.goal.id;
          public          postgres    false    222            �            1259    16505    income    TABLE     �   CREATE TABLE public.income (
    id integer NOT NULL,
    amount numeric,
    date date,
    planned boolean,
    user_id integer
);
    DROP TABLE public.income;
       public         heap    postgres    false            �            1259    16686    income_categories    TABLE     �   CREATE TABLE public.income_categories (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    icon character varying(255) NOT NULL,
    is_fixed boolean NOT NULL,
    user_id integer
);
 %   DROP TABLE public.income_categories;
       public         heap    postgres    false            �            1259    16685    income_categories_id_seq    SEQUENCE     �   CREATE SEQUENCE public.income_categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 /   DROP SEQUENCE public.income_categories_id_seq;
       public          postgres    false    231            u           0    0    income_categories_id_seq    SEQUENCE OWNED BY     U   ALTER SEQUENCE public.income_categories_id_seq OWNED BY public.income_categories.id;
          public          postgres    false    230            �            1259    16504    income_id_seq    SEQUENCE     �   CREATE SEQUENCE public.income_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 $   DROP SEQUENCE public.income_id_seq;
       public          postgres    false    217            v           0    0    income_id_seq    SEQUENCE OWNED BY     ?   ALTER SEQUENCE public.income_id_seq OWNED BY public.income.id;
          public          postgres    false    216            �            1259    16707    investment_categories    TABLE     �   CREATE TABLE public.investment_categories (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    icon character varying(255) NOT NULL,
    is_fixed boolean NOT NULL,
    user_id integer
);
 )   DROP TABLE public.investment_categories;
       public         heap    postgres    false            �            1259    16706    investment_categories_id_seq    SEQUENCE     �   CREATE SEQUENCE public.investment_categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 3   DROP SEQUENCE public.investment_categories_id_seq;
       public          postgres    false    235            w           0    0    investment_categories_id_seq    SEQUENCE OWNED BY     ]   ALTER SEQUENCE public.investment_categories_id_seq OWNED BY public.investment_categories.id;
          public          postgres    false    234            �            1259    16601    sessions    TABLE     ;  CREATE TABLE public.sessions (
    session_id integer NOT NULL,
    username character varying(50) NOT NULL,
    device_id character varying(36) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    last_activity timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    user_id integer
);
    DROP TABLE public.sessions;
       public         heap    postgres    false            �            1259    16600    sessions_session_id_seq    SEQUENCE     �   CREATE SEQUENCE public.sessions_session_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 .   DROP SEQUENCE public.sessions_session_id_seq;
       public          postgres    false    225            x           0    0    sessions_session_id_seq    SEQUENCE OWNED BY     S   ALTER SEQUENCE public.sessions_session_id_seq OWNED BY public.sessions.session_id;
          public          postgres    false    224            �            1259    16695    subscriptions    TABLE     �   CREATE TABLE public.subscriptions (
    id integer NOT NULL,
    user_id integer,
    start_date timestamp without time zone NOT NULL,
    end_date timestamp without time zone NOT NULL,
    is_active boolean NOT NULL
);
 !   DROP TABLE public.subscriptions;
       public         heap    postgres    false            �            1259    16694    subscriptions_id_seq    SEQUENCE     �   CREATE SEQUENCE public.subscriptions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 +   DROP SEQUENCE public.subscriptions_id_seq;
       public          postgres    false    233            y           0    0    subscriptions_id_seq    SEQUENCE OWNED BY     M   ALTER SEQUENCE public.subscriptions_id_seq OWNED BY public.subscriptions.id;
          public          postgres    false    232            �            1259    16630    tracking_state    TABLE     e   CREATE TABLE public.tracking_state (
    id integer NOT NULL,
    state real,
    user_id integer
);
 "   DROP TABLE public.tracking_state;
       public         heap    postgres    false            �            1259    16629    tracking_state_id_seq    SEQUENCE     �   CREATE SEQUENCE public.tracking_state_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 ,   DROP SEQUENCE public.tracking_state_id_seq;
       public          postgres    false    227            z           0    0    tracking_state_id_seq    SEQUENCE OWNED BY     O   ALTER SEQUENCE public.tracking_state_id_seq OWNED BY public.tracking_state.id;
          public          postgres    false    226            �            1259    16496    users    TABLE     �   CREATE TABLE public.users (
    id integer NOT NULL,
    username character varying(50) NOT NULL,
    hashed_password character varying(60) NOT NULL,
    name character varying(50)
);
    DROP TABLE public.users;
       public         heap    postgres    false            �            1259    16495    users_id_seq    SEQUENCE     �   CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 #   DROP SEQUENCE public.users_id_seq;
       public          postgres    false    215            {           0    0    users_id_seq    SEQUENCE OWNED BY     =   ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;
          public          postgres    false    214            �            1259    16533    wealth_fund    TABLE     u   CREATE TABLE public.wealth_fund (
    id integer NOT NULL,
    amount numeric,
    date date,
    user_id integer
);
    DROP TABLE public.wealth_fund;
       public         heap    postgres    false            �            1259    16532    wealth_fund_id_seq    SEQUENCE     �   CREATE SEQUENCE public.wealth_fund_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 )   DROP SEQUENCE public.wealth_fund_id_seq;
       public          postgres    false    221            |           0    0    wealth_fund_id_seq    SEQUENCE OWNED BY     I   ALTER SEQUENCE public.wealth_fund_id_seq OWNED BY public.wealth_fund.id;
          public          postgres    false    220            �           2604    16522 
   expense id    DEFAULT     h   ALTER TABLE ONLY public.expense ALTER COLUMN id SET DEFAULT nextval('public.expense_id_seq'::regclass);
 9   ALTER TABLE public.expense ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    218    219    219            �           2604    16680    expense_categories id    DEFAULT     ~   ALTER TABLE ONLY public.expense_categories ALTER COLUMN id SET DEFAULT nextval('public.expense_categories_id_seq'::regclass);
 D   ALTER TABLE public.expense_categories ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    228    229    229            �           2604    16562    goal id    DEFAULT     b   ALTER TABLE ONLY public.goal ALTER COLUMN id SET DEFAULT nextval('public.goal_id_seq'::regclass);
 6   ALTER TABLE public.goal ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    223    222    223            �           2604    16508 	   income id    DEFAULT     f   ALTER TABLE ONLY public.income ALTER COLUMN id SET DEFAULT nextval('public.income_id_seq'::regclass);
 8   ALTER TABLE public.income ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    216    217    217            �           2604    16689    income_categories id    DEFAULT     |   ALTER TABLE ONLY public.income_categories ALTER COLUMN id SET DEFAULT nextval('public.income_categories_id_seq'::regclass);
 C   ALTER TABLE public.income_categories ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    231    230    231            �           2604    16710    investment_categories id    DEFAULT     �   ALTER TABLE ONLY public.investment_categories ALTER COLUMN id SET DEFAULT nextval('public.investment_categories_id_seq'::regclass);
 G   ALTER TABLE public.investment_categories ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    235    234    235            �           2604    16604    sessions session_id    DEFAULT     z   ALTER TABLE ONLY public.sessions ALTER COLUMN session_id SET DEFAULT nextval('public.sessions_session_id_seq'::regclass);
 B   ALTER TABLE public.sessions ALTER COLUMN session_id DROP DEFAULT;
       public          postgres    false    225    224    225            �           2604    16698    subscriptions id    DEFAULT     t   ALTER TABLE ONLY public.subscriptions ALTER COLUMN id SET DEFAULT nextval('public.subscriptions_id_seq'::regclass);
 ?   ALTER TABLE public.subscriptions ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    233    232    233            �           2604    16633    tracking_state id    DEFAULT     v   ALTER TABLE ONLY public.tracking_state ALTER COLUMN id SET DEFAULT nextval('public.tracking_state_id_seq'::regclass);
 @   ALTER TABLE public.tracking_state ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    227    226    227            �           2604    16499    users id    DEFAULT     d   ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);
 7   ALTER TABLE public.users ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    214    215    215            �           2604    16536    wealth_fund id    DEFAULT     p   ALTER TABLE ONLY public.wealth_fund ALTER COLUMN id SET DEFAULT nextval('public.wealth_fund_id_seq'::regclass);
 =   ALTER TABLE public.wealth_fund ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    221    220    221            [          0    16519    expense 
   TABLE DATA           E   COPY public.expense (id, amount, date, planned, user_id) FROM stdin;
    public          postgres    false    219   pl       e          0    16677    expense_categories 
   TABLE DATA           O   COPY public.expense_categories (id, name, icon, is_fixed, user_id) FROM stdin;
    public          postgres    false    229   �l       _          0    16559    goal 
   TABLE DATA           F   COPY public.goal (id, goal, user_id, need, current_state) FROM stdin;
    public          postgres    false    223   �l       Y          0    16505    income 
   TABLE DATA           D   COPY public.income (id, amount, date, planned, user_id) FROM stdin;
    public          postgres    false    217   m       g          0    16686    income_categories 
   TABLE DATA           N   COPY public.income_categories (id, name, icon, is_fixed, user_id) FROM stdin;
    public          postgres    false    231   >m       k          0    16707    investment_categories 
   TABLE DATA           R   COPY public.investment_categories (id, name, icon, is_fixed, user_id) FROM stdin;
    public          postgres    false    235   lm       a          0    16601    sessions 
   TABLE DATA           g   COPY public.sessions (session_id, username, device_id, created_at, last_activity, user_id) FROM stdin;
    public          postgres    false    225   �m       i          0    16695    subscriptions 
   TABLE DATA           U   COPY public.subscriptions (id, user_id, start_date, end_date, is_active) FROM stdin;
    public          postgres    false    233   �m       c          0    16630    tracking_state 
   TABLE DATA           <   COPY public.tracking_state (id, state, user_id) FROM stdin;
    public          postgres    false    227   )n       W          0    16496    users 
   TABLE DATA           D   COPY public.users (id, username, hashed_password, name) FROM stdin;
    public          postgres    false    215   Fn       ]          0    16533    wealth_fund 
   TABLE DATA           @   COPY public.wealth_fund (id, amount, date, user_id) FROM stdin;
    public          postgres    false    221   #o       }           0    0    expense_categories_id_seq    SEQUENCE SET     G   SELECT pg_catalog.setval('public.expense_categories_id_seq', 2, true);
          public          postgres    false    228            ~           0    0    expense_id_seq    SEQUENCE SET     <   SELECT pg_catalog.setval('public.expense_id_seq', 3, true);
          public          postgres    false    218                       0    0    goal_id_seq    SEQUENCE SET     9   SELECT pg_catalog.setval('public.goal_id_seq', 2, true);
          public          postgres    false    222            �           0    0    income_categories_id_seq    SEQUENCE SET     F   SELECT pg_catalog.setval('public.income_categories_id_seq', 1, true);
          public          postgres    false    230            �           0    0    income_id_seq    SEQUENCE SET     ;   SELECT pg_catalog.setval('public.income_id_seq', 4, true);
          public          postgres    false    216            �           0    0    investment_categories_id_seq    SEQUENCE SET     J   SELECT pg_catalog.setval('public.investment_categories_id_seq', 1, true);
          public          postgres    false    234            �           0    0    sessions_session_id_seq    SEQUENCE SET     F   SELECT pg_catalog.setval('public.sessions_session_id_seq', 59, true);
          public          postgres    false    224            �           0    0    subscriptions_id_seq    SEQUENCE SET     B   SELECT pg_catalog.setval('public.subscriptions_id_seq', 1, true);
          public          postgres    false    232            �           0    0    tracking_state_id_seq    SEQUENCE SET     D   SELECT pg_catalog.setval('public.tracking_state_id_seq', 1, false);
          public          postgres    false    226            �           0    0    users_id_seq    SEQUENCE SET     :   SELECT pg_catalog.setval('public.users_id_seq', 9, true);
          public          postgres    false    214            �           0    0    wealth_fund_id_seq    SEQUENCE SET     @   SELECT pg_catalog.setval('public.wealth_fund_id_seq', 1, true);
          public          postgres    false    220            �           2606    16684 *   expense_categories expense_categories_pkey 
   CONSTRAINT     h   ALTER TABLE ONLY public.expense_categories
    ADD CONSTRAINT expense_categories_pkey PRIMARY KEY (id);
 T   ALTER TABLE ONLY public.expense_categories DROP CONSTRAINT expense_categories_pkey;
       public            postgres    false    229            �           2606    16526    expense expense_pkey 
   CONSTRAINT     R   ALTER TABLE ONLY public.expense
    ADD CONSTRAINT expense_pkey PRIMARY KEY (id);
 >   ALTER TABLE ONLY public.expense DROP CONSTRAINT expense_pkey;
       public            postgres    false    219            �           2606    16564    goal goal_pkey 
   CONSTRAINT     L   ALTER TABLE ONLY public.goal
    ADD CONSTRAINT goal_pkey PRIMARY KEY (id);
 8   ALTER TABLE ONLY public.goal DROP CONSTRAINT goal_pkey;
       public            postgres    false    223            �           2606    16693 (   income_categories income_categories_pkey 
   CONSTRAINT     f   ALTER TABLE ONLY public.income_categories
    ADD CONSTRAINT income_categories_pkey PRIMARY KEY (id);
 R   ALTER TABLE ONLY public.income_categories DROP CONSTRAINT income_categories_pkey;
       public            postgres    false    231            �           2606    16512    income income_pkey 
   CONSTRAINT     P   ALTER TABLE ONLY public.income
    ADD CONSTRAINT income_pkey PRIMARY KEY (id);
 <   ALTER TABLE ONLY public.income DROP CONSTRAINT income_pkey;
       public            postgres    false    217            �           2606    16714 0   investment_categories investment_categories_pkey 
   CONSTRAINT     n   ALTER TABLE ONLY public.investment_categories
    ADD CONSTRAINT investment_categories_pkey PRIMARY KEY (id);
 Z   ALTER TABLE ONLY public.investment_categories DROP CONSTRAINT investment_categories_pkey;
       public            postgres    false    235            �           2606    16608    sessions sessions_pkey 
   CONSTRAINT     \   ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (session_id);
 @   ALTER TABLE ONLY public.sessions DROP CONSTRAINT sessions_pkey;
       public            postgres    false    225            �           2606    16610 (   sessions sessions_username_device_id_key 
   CONSTRAINT     r   ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_username_device_id_key UNIQUE (username, device_id);
 R   ALTER TABLE ONLY public.sessions DROP CONSTRAINT sessions_username_device_id_key;
       public            postgres    false    225    225            �           2606    16700     subscriptions subscriptions_pkey 
   CONSTRAINT     ^   ALTER TABLE ONLY public.subscriptions
    ADD CONSTRAINT subscriptions_pkey PRIMARY KEY (id);
 J   ALTER TABLE ONLY public.subscriptions DROP CONSTRAINT subscriptions_pkey;
       public            postgres    false    233            �           2606    16635 "   tracking_state tracking_state_pkey 
   CONSTRAINT     `   ALTER TABLE ONLY public.tracking_state
    ADD CONSTRAINT tracking_state_pkey PRIMARY KEY (id);
 L   ALTER TABLE ONLY public.tracking_state DROP CONSTRAINT tracking_state_pkey;
       public            postgres    false    227            �           2606    16501    users users_pkey 
   CONSTRAINT     N   ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);
 :   ALTER TABLE ONLY public.users DROP CONSTRAINT users_pkey;
       public            postgres    false    215            �           2606    16503    users users_username_key 
   CONSTRAINT     W   ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);
 B   ALTER TABLE ONLY public.users DROP CONSTRAINT users_username_key;
       public            postgres    false    215            �           2606    16540    wealth_fund wealth_fund_pkey 
   CONSTRAINT     Z   ALTER TABLE ONLY public.wealth_fund
    ADD CONSTRAINT wealth_fund_pkey PRIMARY KEY (id);
 F   ALTER TABLE ONLY public.wealth_fund DROP CONSTRAINT wealth_fund_pkey;
       public            postgres    false    221            �           2606    16720 2   expense_categories expense_categories_user_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.expense_categories
    ADD CONSTRAINT expense_categories_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);
 \   ALTER TABLE ONLY public.expense_categories DROP CONSTRAINT expense_categories_user_id_fkey;
       public          postgres    false    215    229    3237            �           2606    16651    expense expense_user_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.expense
    ADD CONSTRAINT expense_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
 F   ALTER TABLE ONLY public.expense DROP CONSTRAINT expense_user_id_fkey;
       public          postgres    false    219    215    3237            �           2606    16656    goal goal_user_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.goal
    ADD CONSTRAINT goal_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
 @   ALTER TABLE ONLY public.goal DROP CONSTRAINT goal_user_id_fkey;
       public          postgres    false    215    223    3237            �           2606    16715 0   income_categories income_categories_user_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.income_categories
    ADD CONSTRAINT income_categories_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);
 Z   ALTER TABLE ONLY public.income_categories DROP CONSTRAINT income_categories_user_id_fkey;
       public          postgres    false    215    231    3237            �           2606    16646    income income_user_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.income
    ADD CONSTRAINT income_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
 D   ALTER TABLE ONLY public.income DROP CONSTRAINT income_user_id_fkey;
       public          postgres    false    3237    215    217            �           2606    16725 8   investment_categories investment_categories_user_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.investment_categories
    ADD CONSTRAINT investment_categories_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);
 b   ALTER TABLE ONLY public.investment_categories DROP CONSTRAINT investment_categories_user_id_fkey;
       public          postgres    false    215    3237    235            �           2606    16661    sessions sessions_user_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
 H   ALTER TABLE ONLY public.sessions DROP CONSTRAINT sessions_user_id_fkey;
       public          postgres    false    3237    225    215            �           2606    16701 (   subscriptions subscriptions_user_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.subscriptions
    ADD CONSTRAINT subscriptions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);
 R   ALTER TABLE ONLY public.subscriptions DROP CONSTRAINT subscriptions_user_id_fkey;
       public          postgres    false    3237    233    215            �           2606    16666 *   tracking_state tracking_state_user_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.tracking_state
    ADD CONSTRAINT tracking_state_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
 T   ALTER TABLE ONLY public.tracking_state DROP CONSTRAINT tracking_state_user_id_fkey;
       public          postgres    false    227    215    3237            �           2606    16671 $   wealth_fund wealth_fund_user_id_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.wealth_fund
    ADD CONSTRAINT wealth_fund_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
 N   ALTER TABLE ONLY public.wealth_fund DROP CONSTRAINT wealth_fund_user_id_fkey;
       public          postgres    false    3237    215    221            [   #   x�3�4100�4202�54�52�L�4����� ?MP      e   !   x�3�LN,1��L��3�L�4�2B����� Ǜ	&      _   $   x�3��sw�w��4�440�4�3�2�"���� ̨:      Y   &   x�3�4200 Fƺ���F��i�f\&��XDc���� ���      g      x�3�LN,1��L��3�L�4����� A�      k      x�3�LN,1��L��3�L�4����� A�      a   H   x�3���K-/-N-ⴲ2�O.-�ѷ�3�3�4202�54�52W02�2��26�370�43�60& k����� ��      i   '   x�3�4�4202�54�52R00�#���X	W� ��      c      x������ � �      W   �   x�5̱R�0 �9y��85mG�T��J�<�P~�B��	�>}ϡ�7~�¢g��E���XIq_�t��K_��.)P����4�[��l�B���!A�<�o�ý�6�k�8f,��z�)%�vl����O�iW���j־2��@J����{���
����ff�:���b'/?g05�|�����zZ����(�1�B,F�      ]      x������ � �     