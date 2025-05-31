import { Link, useLocation } from '@tanstack/react-router';
import { LuSettings, LuAppWindow } from 'react-icons/lu';
import { cn } from 'src/lib/utils';

interface SettingsLayoutProps {
    children: React.ReactNode;
}

export function SettingsLayout({ children }: SettingsLayoutProps) {
    const location = useLocation();

    const navigationItems = [
        {
            key: 'general',
            label: 'Основные',
            icon: LuSettings,
            href: '/panel/settings',
            isActive: location.pathname === '/panel/settings',
        },
        {
            key: 'external-apps',
            label: 'Внешние приложения',
            icon: LuAppWindow,
            href: '/panel/settings/external-apps',
            isActive: location.pathname === '/panel/settings/external-apps',
        },
    ];

    return (
        <div className="mx-auto w-[800px] pt-8 relative">
            {/* Боковая навигационная панель */}
            <div className="absolute left-[-120px] top-8 w-[100px]">
                <nav className="space-y-2">
                    {navigationItems.map((item) => {
                        const Icon = item.icon;
                        return (
                            <Link
                                key={item.key}
                                to={item.href}

                                className={cn(
                                    'flex flex-col items-center gap-1 p-3 rounded-lg transition-colors text-sm',
                                    'hover:bg-gray-100 hover:text-gray-900',
                                    item.isActive
                                        ? 'bg-blue-50 text-blue-600 border border-blue-200'
                                        : 'text-gray-600'
                                )}
                            >
                                <Icon className="h-5 w-5" />
                                <span className="text-xs text-center leading-tight">
                                    {item.label}
                                </span>
                            </Link>
                        );
                    })}
                </nav>
            </div>

            {/* Основной контент */}
            <div className="w-full">
                {children}
            </div>
        </div>
    );
} 