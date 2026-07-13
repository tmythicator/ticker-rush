import { useAuth } from '@/hooks/useAuth';
import { updateUser } from '@/lib/api';
import { updateUserSchema, type UpdateUserFormData } from '@/lib/schemas';
import { zodResolver } from '@hookform/resolvers/zod';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useForm, useWatch } from 'react-hook-form';

export const useEditProfile = (onClose: () => void) => {
  const { user } = useAuth();
  const queryClient = useQueryClient();

  const {
    register,
    handleSubmit,
    setValue,
    control,
    formState: { errors, dirtyFields },
  } = useForm<UpdateUserFormData>({
    resolver: zodResolver(updateUserSchema),
    defaultValues: {
      firstName: user?.first_name || '',
      lastName: user?.last_name || '',
      website: user?.website || '',
      isPublic: user?.is_public || false,
    },
  });

  const isPublic = useWatch({ control, name: 'isPublic' });

  const mutation = useMutation({
    mutationFn: (data: UpdateUserFormData) => {
      const paths: string[] = [];
      if (dirtyFields.firstName) paths.push('first_name');
      if (dirtyFields.lastName) paths.push('last_name');
      if (dirtyFields.website) paths.push('website');
      if (dirtyFields.isPublic) paths.push('is_public');

      return updateUser({
        first_name: data.firstName,
        last_name: data.lastName,
        website: data.website || '',
        is_public: data.isPublic,
        update_mask: (paths.length > 0 ? { paths } : undefined) as unknown as string[],
      });
    },
    onSuccess: (updatedUser) => {
      queryClient.setQueryData(['user'], updatedUser);
      queryClient.invalidateQueries({ queryKey: ['user'] });
      onClose();
    },
  });

  const onSubmit = handleSubmit((data) => mutation.mutate(data));

  return {
    register,
    onSubmit,
    setValue,
    isPublic,
    errors,
    isPending: mutation.isPending,
    isError: mutation.isError,
    error: mutation.error,
  };
};
