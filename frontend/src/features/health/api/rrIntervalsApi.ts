import { createApi } from '@reduxjs/toolkit/query/react';
import { baseQueryWithReauth } from 'src/store/apiQuery';

export interface RRInterval {
    id: string;
    user_id: string;
    device_id: string;
    rr_interval_ms: number;
    bpm: number;
    created_at: string;
    is_valid: boolean;
}

export interface TimeRange {
    from: string;
    to: string;
}

export interface RRIntervalsResponse {
    intervals: RRInterval[];
    total_count: number;
    valid_count: number;
    time_range: TimeRange;
}

export interface DateRangeParams {
    from: string;    // RFC3339 format
    to: string;      // RFC3339 format
    device_id?: string;
}

export interface CreateBatchRRIntervalsRequest {
    device_id: string;
    intervals: RRIntervalCreateRequest[];
}

export interface RRIntervalCreateRequest {
    rr_interval_ms: number;
    timestamp?: string;
}

export interface CreateBatchRRIntervalsResponse {
    processed_count: number;
    valid_count: number;
    intervals: RRInterval[];
}

export const rrIntervalsApi = createApi({
    reducerPath: 'rrIntervalsApi',
    baseQuery: baseQueryWithReauth,
    tagTypes: ['RRInterval'],
    endpoints: (builder) => ({
        getRRIntervals: builder.query<RRIntervalsResponse, DateRangeParams>({
            query: (params) => ({
                url: '/v1/rr-intervals',
                params,
            }),
            providesTags: ['RRInterval'],
        }),
        createBatchRRIntervals: builder.mutation<CreateBatchRRIntervalsResponse, CreateBatchRRIntervalsRequest>({
            query: (body) => ({
                url: '/v1/rr-intervals/batch',
                method: 'POST',
                body,
            }),
            invalidatesTags: ['RRInterval'],
        }),
    }),
});

export const {
    useGetRRIntervalsQuery,
    useCreateBatchRRIntervalsMutation,
} = rrIntervalsApi; 