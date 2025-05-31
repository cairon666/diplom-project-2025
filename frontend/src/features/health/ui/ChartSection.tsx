import React from 'react';
import { LucideIcon } from 'lucide-react';
import { ChartToggleButton } from './ChartControls';
import { InvalidPeriodAlert, WarningAlert, InfoAlert } from './AlertComponents';
import { ChartVisibilityState } from '../hooks/useChartVisibility';

export interface ChartSectionProps {
    // Основные свойства
    chartName: keyof ChartVisibilityState;
    title: string;
    icon: LucideIcon;
    
    // Состояние видимости и валидности
    isVisible: boolean;
    isValid: boolean;
    onToggle: (chartName: keyof ChartVisibilityState) => void;
    
    // Сообщения
    invalidMessage?: string;
    warningMessage?: React.ReactNode;
    infoMessage?: React.ReactNode;
    
    // Контент графика
    children: React.ReactNode;
    
    // Опциональные настройки
    className?: string;
}

export const ChartSection: React.FC<ChartSectionProps> = ({
    chartName,
    title,
    icon,
    isVisible,
    isValid,
    onToggle,
    invalidMessage,
    warningMessage,
    infoMessage,
    children,
    className = "space-y-4",
}) => {
    return (
        <div className={className}>
            <ChartToggleButton
                chartName={chartName}
                title={title}
                icon={icon}
                isVisible={isVisible}
                isValid={isValid}
                invalidMessage={invalidMessage}
                onToggle={onToggle}
            />
            
            {/* Показываем сообщение о недопустимом периоде */}
            {!isValid && invalidMessage && (
                <InvalidPeriodAlert title={`Недопустимый период для ${title.toLowerCase()}`}>
                    {invalidMessage}
                </InvalidPeriodAlert>
            )}
            
            {/* Контент графика показываем только если он видим и валиден */}
            {isVisible && isValid && (
                <>
                    {/* Предупреждения */}
                    {warningMessage && (
                        <WarningAlert>
                            {warningMessage}
                        </WarningAlert>
                    )}
                    
                    {/* Информационные сообщения */}
                    {infoMessage && (
                        <InfoAlert>
                            {infoMessage}
                        </InfoAlert>
                    )}
                    
                    {/* Сам график */}
                    {children}
                </>
            )}
        </div>
    );
}; 