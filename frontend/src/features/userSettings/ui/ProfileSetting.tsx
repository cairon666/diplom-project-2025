import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { setPartialCurrentUser, currentUserApi } from 'src/entities/user';
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

import { userSettingApi, useUpdateProfileSettingMutation, useSettingQuery } from '../api';

const ProfileSettingFormSchema = z.object({
    firstName: z.string().min(1, 'Имя обязательно'),
    lastName: z.string().min(1, 'Фамилия обязательна'),
});

export function ProfileSetting() {
    const dispatch = useAppDispatch();
    const { data } = useSettingQuery();
    const [updateProfile, { isLoading, error }] = useUpdateProfileSettingMutation();

    const form = useForm<z.infer<typeof ProfileSettingFormSchema>>({
        resolver: zodResolver(ProfileSettingFormSchema),
        defaultValues: {
            firstName: data?.first_name || '',
            lastName: data?.last_name || '',
        },
    });

    const onUpdateProfile = (firstName: string, lastName: string) => {
        updateProfile({
            first_name: firstName,
            last_name: lastName,
        })
            .unwrap()
            .then(() => {
                dispatch(
                    userSettingApi.util.updateQueryData(
                        'setting',
                        undefined,
                        (draft) => {
                            draft.first_name = firstName;
                            draft.last_name = lastName;
                        },
                    )
                );
                dispatch(
                    currentUserApi.util.updateQueryData('getUser', undefined, (draft) => {
                        draft.first_name = firstName;
                        draft.last_name = lastName;
                    })
                );
                dispatch(
                    setPartialCurrentUser({
                        first_name: firstName,
                        last_name: lastName,
                    }),
                );
                toast.success('Профиль успешно обновлен');
            })
            .catch(() => {
                toast.error('Произошла ошибка при обновлении профиля');
            });
    };

    const onSubmit = form.handleSubmit((data) =>
        onUpdateProfile(data.firstName, data.lastName),
    );

    if (!data) {
        return null;
    }

    return (
        <div className="flex flex-col gap-4 w-full">
            <div className="flex flex-col gap-4">
                <div className="flex flex-col">
                    <h3 className="text-lg font-medium">Профиль</h3>
                    <p className="text-sm text-muted-foreground">
                        Изменить основную информацию профиля
                    </p>
                </div>
                
                <Form {...form}>
                    <form onSubmit={onSubmit} className="flex flex-col gap-4">
                        <div className="flex gap-4">
                            <FormField
                                control={form.control}
                                name="firstName"
                                render={({ field }) => (
                                    <FormItem className="flex flex-col gap-1 flex-1">
                                        <FormLabel>Имя</FormLabel>
                                        <FormControl>
                                            <Input {...field} placeholder="Введите имя" />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name="lastName"
                                render={({ field }) => (
                                    <FormItem className="flex flex-col gap-1 flex-1">
                                        <FormLabel>Фамилия</FormLabel>
                                        <FormControl>
                                            <Input
                                                {...field}
                                                placeholder="Введите фамилию"
                                            />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                        </div>
                        <div className="flex justify-end">
                            <Button type="submit" disabled={isLoading} className="w-fit">
                                Обновить профиль
                            </Button>
                        </div>
                    </form>
                </Form>
            </div>
        </div>
    );
}
