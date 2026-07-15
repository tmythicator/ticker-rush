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
import styles from './EditProfileModal.module.css';

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

          <ModalBody className={styles.formBody}>
            <div className={styles.section}>
              <h3 className={styles.sectionTitle}>
                Personal Information
              </h3>
              <div className={styles.nameRow}>
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

            <div className={styles.section}>
              <h3 className={styles.sectionTitle}>
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
              className={styles.footerButton}
              data-testid="edit-profile-cancel"
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={isPending}
              className={styles.footerButton}
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
