import Image from 'next/image';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { useContext } from 'react';
import GlobalAppContext from '../layout/context/GlobalAppContext';

function Navbar() {
  const { displayBanner } = useContext(GlobalAppContext);
  const router = useRouter();
  const nav = [
    { label: 'Dashboard', href: '/dashboard' },
    { label: 'Inventory', href: '/inventory' },
    { label: 'Code', href: '/code' },
    { label: 'Operations', href: '/operations' },
    { label: 'SaaS', href: '/saas' },
    { label: 'Security', href: '/security' },
    { label: 'Map', href: '/map' }
  ];
  return (
    <nav
      className={`fixed ${
        displayBanner ? 'top-[72px]' : 'top-0'
      } z-30 flex w-full items-center justify-between gap-10 border-b border-black-200/30 bg-white px-6 py-4 xl:pr-8 2xl:pr-24`}
    >
      <div className="flex items-center gap-8 text-sm font-semibold text-black-400">
        <Link href="/dashboard">
          <Image
            src="/assets/img/komiser.svg"
            width={40}
            height={40}
            alt="Komiser logo"
          />
        </Link>
        {nav.map((navItem, idx) => (
          <Link
            key={idx}
            href={navItem.href}
            className={
              router.pathname === navItem.href
                ? 'text-primary'
                : 'text-black-400'
            }
          >
            {navItem.label}
          </Link>
        ))}
      </div>
      <div className="flex gap-4 text-sm font-medium text-black-900 lg:gap-10">
        <a
          className="hidden items-center gap-2 transition-colors hover:text-primary md:flex"
          href="https://www.tailwarden.com/changelog?utm_source=komiser&utm_medium=referral&utm_campaign=static"
          target="_blank"
          rel="noopener noreferrer"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth="1.5"
            stroke="currentColor"
            className="h-5 w-5"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M10.343 3.94c.09-.542.56-.94 1.11-.94h1.093c.55 0 1.02.398 1.11.94l.149.894c.07.424.384.764.78.93.398.164.855.142 1.205-.108l.737-.527a1.125 1.125 0 011.45.12l.773.774c.39.389.44 1.002.12 1.45l-.527.737c-.25.35-.272.806-.107 1.204.165.397.505.71.93.78l.893.15c.543.09.94.56.94 1.109v1.094c0 .55-.397 1.02-.94 1.11l-.893.149c-.425.07-.765.383-.93.78-.165.398-.143.854.107 1.204l.527.738c.32.447.269 1.06-.12 1.45l-.774.773a1.125 1.125 0 01-1.449.12l-.738-.527c-.35-.25-.806-.272-1.203-.107-.397.165-.71.505-.781.929l-.149.894c-.09.542-.56.94-1.11.94h-1.094c-.55 0-1.019-.398-1.11-.94l-.148-.894c-.071-.424-.384-.764-.781-.93-.398-.164-.854-.142-1.204.108l-.738.527c-.447.32-1.06.269-1.45-.12l-.773-.774a1.125 1.125 0 01-.12-1.45l.527-.737c.25-.35.273-.806.108-1.204-.165-.397-.505-.71-.93-.78l-.894-.15c-.542-.09-.94-.56-.94-1.109v-1.094c0-.55.398-1.02.94-1.11l.894-.149c.424-.07.765-.383.93-.78.165-.398.143-.854-.107-1.204l-.527-.738a1.125 1.125 0 01.12-1.45l.773-.773a1.125 1.125 0 011.45-.12l.737.527c.35.25.807.272 1.204.107.397-.165.71-.505.78-.929l.15-.894z"
            />
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
            />
          </svg>
          Settings
        </a>
        <a
          className="flex items-center gap-2 rounded-lg bg-[#5865F2] px-4 py-2 text-white transition-colors hover:bg-[#4f5be2]"
          href="https://discord.tailwarden.com"
          target="_blank"
          rel="noopener noreferrer"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            fill="none"
            viewBox="0 0 24 24"
          >
            <path
              fill="currentColor"
              d="M18.93 4.935a16.457 16.457 0 00-4.07-1.266.062.062 0 00-.066.031c-.175.314-.37.723-.506 1.044a15.183 15.183 0 00-4.573 0c-.136-.328-.338-.73-.515-1.044a.064.064 0 00-.065-.031 16.413 16.413 0 00-4.07 1.266.058.058 0 00-.028.023c-2.593 3.885-3.303 7.674-2.954 11.417a.069.069 0 00.026.047 16.565 16.565 0 004.994 2.531.065.065 0 00.07-.023c.385-.527.728-1.082 1.022-1.666a.064.064 0 00-.035-.089 10.906 10.906 0 01-1.56-.745.064.064 0 01-.007-.107c.105-.079.21-.16.31-.244a.061.061 0 01.065-.008c3.273 1.498 6.817 1.498 10.051 0a.062.062 0 01.066.008c.1.082.204.165.31.244a.064.064 0 01-.005.107c-.499.292-1.017.538-1.561.744a.064.064 0 00-.034.09c.3.583.643 1.139 1.02 1.666a.063.063 0 00.07.023 16.51 16.51 0 005.003-2.531.065.065 0 00.026-.047c.417-4.326-.699-8.084-2.957-11.416a.05.05 0 00-.026-.024zM8.684 14.096c-.985 0-1.797-.907-1.797-2.022 0-1.114.796-2.021 1.797-2.021 1.01 0 1.813.915 1.798 2.021 0 1.115-.796 2.022-1.798 2.022zm6.646 0c-.986 0-1.797-.907-1.797-2.022 0-1.114.796-2.021 1.797-2.021 1.009 0 1.813.915 1.797 2.021 0 1.115-.788 2.022-1.797 2.022z"
            ></path>
          </svg>
          Community (Discord)
        </a>
      </div>
    </nav>
  );
}

export default Navbar;
