import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { useAppDispatch } from 'src/store';
import { Button } from 'ui/button';
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from 'ui/form';
import { Input } from 'ui/input';
import { z } from 'zod';

import { 
    userSettingApi, 
    useSettingQuery, 
    useUpdatePasswordSettingMutation 
} from '../api';

const PasswordSettingFormSchema = z.object({
    password: z.string().min(6, 'Пароль должен содержать минимум 6 символов'),
    confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
    message: "Пароли не совпадают",
    path: ["confirmPassword"],
});

export function PasswordSetting() {
    const dispatch = useAppDispatch();
    const { data } = useSettingQuery();
    const [updatePassword, { isLoading }] = useUpdatePasswordSettingMutation();

    const form = useForm<z.infer<typeof PasswordSettingFormSchema>>({
        resolver: zodResolver(PasswordSettingFormSchema),
        defaultValues: {
            password: '',
            confirmPassword: '',
        },
    });

    const onUpdatePassword = (password: string) => {
        updatePassword({ password })
            .unwrap()
            .then(() => {
                dispatch(
                    userSettingApi.util.updateQueryData(
                        'setting',
                        undefined,
                        (draft) => {
                            draft.has_password = true;
                        },
                    )
                );
                form.reset();
                toast.success(data?.has_password ? 'Пароль успешно изменен' : 'Пароль успешно установлен');
            })
            .catch(() => {
                toast.error('Произошла ошибка при изменении пароля');
            });
    };

    const onSubmit = form.handleSubmit((data) => onUpdatePassword(data.password));

    if (!data) {
        return null;
    }

    return (    
        <div className="flex flex-col gap-4 w-full">
            <div className="flex flex-col gap-4">
                <div className="flex flex-col">
                    <h3 className="text-lg font-medium">Пароль</h3>
                    <p className="text-sm text-muted-foreground">
                        {data.has_password 
                            ? 'Изменить пароль для входа в аккаунт' 
                            : 'Установить пароль для входа в аккаунт'
                        }
                    </p>
                </div>
                
                <Form {...form}>
                    <form onSubmit={onSubmit} className="flex flex-col gap-4">
                        <FormField
                            control={form.control}
                            name="password"
                            render={({ field }) => (
                                <FormItem className="flex flex-col gap-1">
                                    <FormLabel>
                                        {data.has_password ? 'Новый пароль' : 'Пароль'}
                                    </FormLabel>
                                    <FormControl>
                                        <Input 
                                            {...field} 
                                            type="password"
                                            placeholder="Введите пароль" 
                                        />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <FormField
                            control={form.control}
                            name="confirmPassword"
                            render={({ field }) => (
                                <FormItem className="flex flex-col gap-1">
                                    <FormLabel>Подтвердить пароль</FormLabel>
                                    <FormControl>
                                        <Input 
                                            {...field} 
                                            type="password"
                                            placeholder="Подтвердите пароль" 
                                        />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <div className="flex justify-end">
                            <Button type="submit" disabled={isLoading} className="w-fit">
                                {isLoading 
                                    ? (data.has_password ? 'Изменяем...' : 'Устанавливаем...') 
                                    : (data.has_password ? 'Изменить пароль' : 'Установить пароль')
                                }
                            </Button>
                        </div>
                    </form>
                </Form>
            </div>
        </div>
    );
}
