import { createApi } from '@reduxjs/toolkit/query/react';
import { baseQueryWithReauth } from 'src/store/apiQuery';

export interface DeviceListItem {
    id: string;
    user_id: string;
    device_name: string;
    created_at: string;
}

export interface GetDeviceListResponse {
    devices: DeviceListItem[];
}

export interface DeleteDeviceRequest {
    id: string;
}

export const deviceApi = createApi({
    reducerPath: 'deviceApi',
    baseQuery: baseQueryWithReauth,
    tagTypes: ['Device'],
    endpoints: (builder) => ({
        getDeviceList: builder.query<GetDeviceListResponse, void>({
            query: () => ({
                url: `/v1/user/devices`,
                method: 'GET',
            }),
            providesTags: ['Device'],
        }),
        deleteDevice: builder.mutation<unknown, DeleteDeviceRequest>({
            query: ({ id }) => ({
                url: `/v1/user/devices/${id}`,
                method: 'DELETE',
            }),
            invalidatesTags: ['Device'],
        }),
    }),
});

export const {
    useGetDeviceListQuery,
    useDeleteDeviceMutation,
} = deviceApi; 