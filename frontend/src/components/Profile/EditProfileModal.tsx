import { Modal } from '@/components/Modal';
import { FormField } from '@/components/ui/form-field';
import { useAuth } from '@/hooks/useAuth';
import { updateUser } from '@/lib/api';
import { updateUserSchema, type UpdateUserFormData } from '@/lib/schemas';
import { zodResolver } from '@hookform/resolvers/zod';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useForm, useWatch } from 'react-hook-form';

interface EditProfileModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const EditProfileModal = ({ isOpen, onClose }: EditProfileModalProps) => {
  const { user } = useAuth();
  const queryClient = useQueryClient();

  const {
    register,
    handleSubmit,
    setValue,
    control,
    formState: { errors },
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
    mutationFn: (data: UpdateUserFormData) =>
      updateUser({
        first_name: data.firstName,
        last_name: data.lastName,
        website: data.website || '',
        is_public: data.isPublic,
      }),
    onSuccess: (updatedUser) => {
      queryClient.setQueryData(['user'], updatedUser);
      queryClient.invalidateQueries({ queryKey: ['user'] });
      onClose();
    },
  });

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Edit Profile">
      <form onSubmit={handleSubmit((data) => mutation.mutate(data))} className="space-y-6">
        <div className="space-y-4">
          <h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wider">
            Personal Information
          </h3>
          <div className="grid grid-cols-2 gap-4">
            <FormField
              label="First Name"
              id="firstName"
              register={register}
              error={errors.firstName?.message}
              placeholder="John"
            />
            <FormField
              label="Last Name"
              id="lastName"
              register={register}
              error={errors.lastName?.message}
              placeholder="Doe"
            />
          </div>
        </div>

        <div className="space-y-4">
          <h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wider">
            Public Profile
          </h3>
          <FormField
            label="Website"
            id="website"
            register={register}
            error={errors.website?.message}
            placeholder="https://example.com"
          />

          <div
            className="bg-muted/30 p-4 rounded-xl border border-border/50 flex items-center justify-between group cursor-pointer hover:bg-muted/50 transition-colors"
            onClick={() => setValue('isPublic', !isPublic, { shouldDirty: true })}
          >
            <div className="space-y-1">
              <div className="flex items-center gap-2">
                <span className="font-medium text-foreground">Profile Visibility</span>
                <span className="text-xs px-2 py-0.5 rounded-full bg-primary/10 text-primary font-medium">
                  {isPublic ? 'Public' : 'Private'}
                </span>
              </div>
              <p className="text-xs text-muted-foreground">
                When public, your portfolio allocation is visible on the leaderboard.
              </p>
            </div>
            <input
              type="checkbox"
              {...register('isPublic')}
              className="h-5 w-5 rounded border-border text-primary focus:ring-primary/50 bg-background accent-primary cursor-pointer"
            />
          </div>
        </div>

        {mutation.isError && (
          <div className="p-3 rounded-lg bg-destructive/10 border border-destructive/20 text-destructive text-sm font-medium">
            {mutation.error instanceof Error ? mutation.error.message : 'Failed to update profile'}
          </div>
        )}

        <div className="flex gap-3 pt-4">
          <button
            type="button"
            onClick={onClose}
            className="flex-1 px-4 py-2.5 bg-muted text-foreground font-medium rounded-xl hover:bg-muted/80 transition-colors"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={mutation.isPending}
            className="flex-1 px-4 py-2.5 bg-primary text-primary-foreground font-medium rounded-xl hover:bg-primary/90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed shadow-lg shadow-primary/20"
          >
            {mutation.isPending ? 'Saving...' : 'Save Changes'}
          </button>
        </div>
      </form>
    </Modal>
  );
};
