import React from 'react';
import { Alert, AlertDescription } from 'ui/alert';
import { Card, CardContent } from 'ui/card';
import { AlertTriangle, Info } from 'lucide-react';

export interface AlertProps {
    children: React.ReactNode;
}

export const WarningAlert: React.FC<AlertProps> = ({ children }) => (
    <Alert className="border-amber-200 bg-amber-50">
        <AlertTriangle className="h-4 w-4 text-amber-600" />
        <AlertDescription className="text-amber-800">
            {children}
        </AlertDescription>
    </Alert>
);

export const InfoAlert: React.FC<AlertProps> = ({ children }) => (
    <Alert className="border-blue-200 bg-blue-50">
        <Info className="h-4 w-4 text-blue-600" />
        <AlertDescription className="text-blue-800">
            {children}
        </AlertDescription>
    </Alert>
);

export interface InvalidPeriodAlertProps {
    title: string;
    children: React.ReactNode;
}

export const InvalidPeriodAlert: React.FC<InvalidPeriodAlertProps> = ({ title, children }) => (
    <Card>
        <CardContent className="py-6">
            <Alert variant="destructive">
                <AlertTriangle className="h-4 w-4" />
                <AlertDescription>
                    <div>
                        <strong>{title}</strong>
                        <br />
                        {children}
                    </div>
                </AlertDescription>
            </Alert>
        </CardContent>
    </Card>
);

export interface ErrorAlertProps {
    title?: string;
    error: unknown;
    getErrorMessage: (error: unknown) => string;
}

export const ErrorAlert: React.FC<ErrorAlertProps> = ({ 
    title = "Ошибка при загрузке данных", 
    error, 
    getErrorMessage 
}) => (
    <Alert variant="destructive">
        <AlertDescription>
            {title}: {getErrorMessage(error)}
        </AlertDescription>
    </Alert>
); 