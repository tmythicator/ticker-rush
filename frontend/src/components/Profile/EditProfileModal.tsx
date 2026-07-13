import {
  Modal,
  ModalCard,
  ModalHeader,
  ModalTitle,
  ModalCloseButton,
  ModalBody,
  ModalFooter,
} from '@/components/Modal';
import { Button } from '@/components/shared/Button';
import { FormInput } from '@/components/shared/FormInput';
import { ErrorMessage } from '@/components/shared/ErrorMessage';
import { useEditProfile } from './useEditProfile';
import { ProfileVisibilityToggle } from './ProfileVisibilityToggle';
import { DeleteProfileSection } from './DeleteProfileSection';

interface EditProfileModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const EditProfileModal = ({ isOpen, onClose }: EditProfileModalProps) => {
  const { register, onSubmit, setValue, isPublic, errors, isPending, isError, error } =
    useEditProfile(onClose);

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalCard size="md">
        <form onSubmit={onSubmit}>
          <ModalHeader>
            <ModalTitle>Edit Profile</ModalTitle>
            <ModalCloseButton />
          </ModalHeader>

          <ModalBody className="space-y-6">
            <div className="space-y-4">
              <h3 className="text-sm font-semibold uppercase tracking-wider text-muted-foreground">
                Personal Information
              </h3>
              <div className="grid grid-cols-2 gap-4">
                <FormInput
                  label="First Name"
                  id="firstName"
                  placeholder="John"
                  register={register}
                  error={errors.firstName?.message}
                  data-testid="first-name-input"
                />
                <FormInput
                  label="Last Name"
                  id="lastName"
                  placeholder="Doe"
                  register={register}
                  error={errors.lastName?.message}
                  data-testid="last-name-input"
                />
              </div>
            </div>

            <div className="space-y-4">
              <h3 className="text-sm font-semibold uppercase tracking-wider text-muted-foreground">
                Public Profile
              </h3>
              <FormInput
                label="Website"
                id="website"
                placeholder="https://example.com"
                register={register}
                error={errors.website?.message}
                data-testid="website-input"
              />

              <ProfileVisibilityToggle
                isPublic={isPublic}
                onToggle={() => setValue('isPublic', !isPublic, { shouldDirty: true })}
                checkboxProps={register('isPublic')}
              />
            </div>

            <DeleteProfileSection onSuccess={onClose} />

            {isError && (
              <ErrorMessage
                message={error instanceof Error ? error.message : 'Failed to update profile'}
              />
            )}
          </ModalBody>

          <ModalFooter>
            <Button
              type="button"
              variant="secondary"
              onClick={onClose}
              className="flex-1"
              data-testid="edit-profile-cancel"
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={isPending}
              className="flex-1"
              data-testid="edit-profile-submit"
            >
              {isPending ? 'Saving...' : 'Save Changes'}
            </Button>
          </ModalFooter>
        </form>
      </ModalCard>
    </Modal>
  );
};
