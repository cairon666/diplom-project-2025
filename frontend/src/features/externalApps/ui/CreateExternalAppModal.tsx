import { useState } from 'react';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { LuCopy, LuCheck } from 'react-icons/lu';
import { Button } from 'ui/button';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from 'ui/dialog';
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from 'ui/form';
import { Input } from 'ui/input';
import { Checkbox } from 'ui/checkbox';
import { z } from 'zod';

import { useCreateExternalAppMutation, CreateExternalAppResponse } from '../api';

const CreateExternalAppFormSchema = z.object({
    name: z.string().min(1, 'Название обязательно'),
    roles: z.array(z.string()).min(1, 'Выберите хотя бы одну роль'),
});

interface CreateExternalAppModalProps {
    children: React.ReactNode;
}

export function CreateExternalAppModal({ children }: CreateExternalAppModalProps) {
    const [open, setOpen] = useState(false);
    const [createdApp, setCreatedApp] = useState<CreateExternalAppResponse | null>(null);
    const [copied, setCopied] = useState(false);
    const [createExternalApp, { isLoading }] = useCreateExternalAppMutation();

    const form = useForm<z.infer<typeof CreateExternalAppFormSchema>>({
        resolver: zodResolver(CreateExternalAppFormSchema),
        defaultValues: {
            name: '',
            roles: [],
        },
    });

    const roles = [
        { id: 'external_app_reader', label: 'Чтение данных', description: 'Доступ к просмотру данных' },
        { id: 'external_app_writer', label: 'Запись данных', description: 'Доступ к изменению данных' },
    ];

    const onSubmit = form.handleSubmit(async (data) => {
        try {
            const response = await createExternalApp(data).unwrap();
            setCreatedApp(response);
            toast.success('Внешнее приложение успешно создано');
        } catch (error) {
            toast.error('Произошла ошибка при создании приложения');
        }
    });

    const copyApiKey = async () => {
        if (createdApp) {
            await navigator.clipboard.writeText(createdApp.api_key);
            setCopied(true);
            toast.success('API ключ скопирован в буфер обмена');
            setTimeout(() => setCopied(false), 2000);
        }
    };

    const handleClose = () => {
        setOpen(false);
        setCreatedApp(null);
        setCopied(false);
        form.reset();
    };

    if (createdApp) {
        return (
            <Dialog open={open} onOpenChange={handleClose}>
                <DialogTrigger asChild>{children}</DialogTrigger>
                <DialogContent className="bg-white border border-gray-200 shadow-xl max-w-md">
                    <DialogHeader>
                        <DialogTitle className="text-gray-900">Приложение создано!</DialogTitle>
                        <DialogDescription className="text-gray-600">
                            Сохраните API ключ в безопасном месте. Он больше не будет показан.
                        </DialogDescription>
                    </DialogHeader>
                    
                    <div className="space-y-4">
                        <div>
                            <label className="text-sm font-medium text-gray-900">API ключ</label>
                            <div className="flex items-center gap-2 mt-1">
                                <Input 
                                    value={createdApp.api_key} 
                                    readOnly 
                                    className="font-mono text-sm"
                                />
                                <Button
                                    size="icon"
                                    variant="outline"
                                    onClick={copyApiKey}
                                    className="shrink-0"
                                >
                                    {copied ? (
                                        <LuCheck className="h-4 w-4 text-green-600" />
                                    ) : (
                                        <LuCopy className="h-4 w-4" />
                                    )}
                                </Button>
                            </div>
                        </div>
                    </div>

                    <DialogFooter>
                        <Button onClick={handleClose}>Готово</Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        );
    }

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>{children}</DialogTrigger>
            <DialogContent className="bg-white border border-gray-200 shadow-xl">
                <DialogHeader>
                    <DialogTitle className="text-gray-900">Создать внешнее приложение</DialogTitle>
                    <DialogDescription className="text-gray-600">
                        Создайте новое внешнее приложение для интеграции с API.
                    </DialogDescription>
                </DialogHeader>

                <Form {...form}>
                    <form onSubmit={onSubmit} className="space-y-4">
                        <FormField
                            control={form.control}
                            name="name"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Название приложения</FormLabel>
                                    <FormControl>
                                        <Input {...field} placeholder="Мое приложение" />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />

                        <FormField
                            control={form.control}
                            name="roles"
                            render={() => (
                                <FormItem>
                                    <FormLabel>Разрешения</FormLabel>
                                    <div className="space-y-3">
                                        {roles.map((role) => (
                                            <FormField
                                                key={role.id}
                                                control={form.control}
                                                name="roles"
                                                render={({ field }) => (
                                                    <FormItem className="flex flex-row items-start space-x-3 space-y-0">
                                                        <FormControl>
                                                            <Checkbox
                                                                checked={field.value?.includes(role.id)}
                                                                onCheckedChange={(checked) => {
                                                                    return checked
                                                                        ? field.onChange([...field.value, role.id])
                                                                        : field.onChange(
                                                                              field.value?.filter(
                                                                                  (value) => value !== role.id
                                                                              )
                                                                          );
                                                                }}
                                                            />
                                                        </FormControl>
                                                        <div className="space-y-1 leading-none">
                                                            <FormLabel className="font-normal">
                                                                {role.label}
                                                            </FormLabel>
                                                            <p className="text-sm text-muted-foreground">
                                                                {role.description}
                                                            </p>
                                                        </div>
                                                    </FormItem>
                                                )}
                                            />
                                        ))}
                                    </div>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />

                        <DialogFooter>
                            <Button type="button" variant="outline" onClick={() => setOpen(false)}>
                                Отмена
                            </Button>
                            <Button type="submit" disabled={isLoading}>
                                {isLoading ? 'Создаем...' : 'Создать'}
                            </Button>
                        </DialogFooter>
                    </form>
                </Form>
            </DialogContent>
        </Dialog>
    );
} 