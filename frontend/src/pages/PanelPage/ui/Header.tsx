import { Link, useLocation, useNavigate } from '@tanstack/react-router';
import { LuSettings, LuLayoutPanelLeft, LuSmartphone, LuLogOut, LuActivity } from 'react-icons/lu';
import { selectCurrentUser, authSlice } from 'src/entities/user';
import { cn } from 'src/lib/utils';
import { useAppSelector, useAppDispatch } from 'src/store';
import { Avatar, AvatarFallback } from 'ui/avatar';
import { Button } from 'ui/button';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from 'ui/dropdown-menu';
import {
    NavigationMenu,
    NavigationMenuIndicator,
    NavigationMenuItem,
    NavigationMenuLink,
    NavigationMenuList,
} from 'ui/navigation-menu';

export function Header() {
    const location = useLocation();
    const navigate = useNavigate();
    const dispatch = useAppDispatch();
    const currentUser = useAppSelector(selectCurrentUser);

    // Функция для сокращения длинных текстов
    const truncateText = (text: string, maxLength: number) => {
        if (text.length <= maxLength) return text;
        return text.substring(0, maxLength) + '...';
    };

    const fullName = `${currentUser.first_name} ${currentUser.last_name}`;
    const shortName = truncateText(fullName, 20);
    const shortEmail = truncateText(currentUser.email, 25);

    // Обработчик выхода из системы
    const handleLogout = () => {
        // Очищаем состояние в Redux (включая localStorage)
        dispatch(authSlice.actions.logout());
        
        // Перенаправляем на страницу аутентификации
        navigate({ to: '/auth' });
    };

    // Обработчик перехода в настройки
    const handleSettingsClick = () => {
        navigate({ to: '/panel/settings' });
    };

    const navigationItems = [
        {
            href: '/panel',
            icon: LuLayoutPanelLeft,
            label: 'Панель',
            isActive: location.pathname === '/panel',
        },
        {
            href: '/panel/comparative-analysis',
            icon: LuActivity,
            label: 'Сравнительный анализ',
            isActive: location.pathname === '/panel/comparative-analysis',
        },
        {
            href: '/panel/devices',
            icon: LuSmartphone,
            label: 'Устройства',
            isActive: location.pathname === '/panel/devices',
        },
        {
            href: '/panel/settings',
            icon: LuSettings,
            label: 'Настройки',
            isActive: location.pathname.startsWith('/panel/settings'),
        },
    ];

    return (
        <header className="bg-gradient-to-r from-slate-900 via-slate-800 to-slate-900 backdrop-blur-sm border-b border-slate-700/50 h-16 relative w-full sticky top-0 z-50 shadow-lg">
            <div className="container mx-auto h-[64px] px-4">
                <div className="flex items-center justify-between h-full">
                    {/* Spacer для выравнивания навигации по центру */}
                    <div className="flex-1"></div>

                    {/* Navigation - по центру */}
                    <nav className="hidden md:flex items-center space-x-1">
                        {navigationItems.map((item) => {
                            const Icon = item.icon;
                            return (
                                <Link
                                    key={item.href}
                                    to={item.href}
                                    className={cn(
                                        'flex items-center space-x-2 px-4 py-2 rounded-lg transition-all duration-200 text-sm font-medium',
                                        item.isActive
                                            ? 'bg-blue-600 text-white shadow-lg shadow-blue-600/25'
                                            : 'text-slate-300 hover:text-white hover:bg-slate-700/50'
                                    )}
                                >
                                    <Icon className="w-4 h-4" />
                                    <span>{item.label}</span>
                                </Link>
                            );
                        })}
                    </nav>

                    {/* User Menu */}
                    <div className="flex items-center flex-1 justify-end">
                        {/* User Menu */}
                        <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                                <Button
                                    variant="ghost"
                                    className="flex items-center space-x-3 px-3 py-2 hover:bg-slate-700/50 text-slate-300 hover:text-white"
                                >
                                    <Avatar className="w-8 h-8">
                                        <AvatarFallback className="bg-gradient-to-br from-blue-500 to-purple-600 text-white text-sm font-semibold">
                                            {currentUser.first_name[0]}
                                            {currentUser.last_name[0]}
                                        </AvatarFallback>
                                    </Avatar>
                                    <div className="hidden sm:block text-left max-w-[180px]">
                                        <div className="text-sm font-medium truncate">
                                            {shortName}
                                        </div>
                                        <div className="text-xs text-slate-400 truncate">
                                            {shortEmail}
                                        </div>
                                    </div>
                                </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent className="w-56 bg-white border border-slate-200 shadow-xl" align="end">
                                <div className="flex items-center space-x-3 p-3">
                                    <Avatar className="w-10 h-10">
                                        <AvatarFallback className="bg-gradient-to-br from-blue-500 to-purple-600 text-white font-semibold">
                                            {currentUser.first_name[0]}
                                            {currentUser.last_name[0]}
                                        </AvatarFallback>
                                    </Avatar>
                                    <div className="min-w-0 flex-1">
                                        <div className="font-medium text-slate-900 truncate">
                                            {fullName}
                                        </div>
                                        <div className="text-sm text-slate-500 truncate">
                                            {currentUser.email}
                                        </div>
                                    </div>
                                </div>
                                <DropdownMenuSeparator />
                                <DropdownMenuItem 
                                    className="text-slate-700 hover:bg-slate-50 cursor-pointer"
                                    onClick={handleSettingsClick}
                                >
                                    <LuSettings className="w-4 h-4 mr-2" />
                                    Настройки профиля
                                </DropdownMenuItem>
                                <DropdownMenuSeparator />
                                <DropdownMenuItem 
                                    className="text-red-600 hover:bg-red-50 cursor-pointer"
                                    onClick={handleLogout}
                                >
                                    <LuLogOut className="w-4 h-4 mr-2" />
                                    Выйти
                                </DropdownMenuItem>
                            </DropdownMenuContent>
                        </DropdownMenu>
                    </div>
                </div>
            </div>
        </header>
    );
}
