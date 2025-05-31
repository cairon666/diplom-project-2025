import React from 'react';
import { LoadingSpinner } from './loading-spinner';

interface LoadingOverlayProps {
    isLoading?: boolean;
    children: React.ReactNode;
    className?: string;
    loadingText?: string;
    overlay?: boolean;
}

export function LoadingOverlay({ 
    isLoading = false, 
    children, 
    className = '',
    loadingText = 'Загрузка...',
    overlay = true
}: LoadingOverlayProps) {
    return (
        <div className={`relative ${className}`}>
            {children}
            
            {isLoading && overlay && (
                <div className="absolute inset-0 bg-white/70 backdrop-blur-[2px] flex items-center justify-center z-20 rounded-lg">
                    <div className="flex flex-col items-center gap-2 bg-white/90 px-4 py-3 rounded-lg shadow-sm border border-gray-200">
                        <LoadingSpinner />
                        <span className="text-xs font-medium text-gray-600">
                            {loadingText}
                        </span>
                    </div>
                </div>
            )}
            
            {isLoading && !overlay && (
                <div className="flex items-center justify-center py-8">
                    <div className="flex flex-col items-center gap-2">
                        <LoadingSpinner />
                        <span className="text-sm font-medium text-gray-600">
                            {loadingText}
                        </span>
                    </div>
                </div>
            )}
        </div>
    );
} 