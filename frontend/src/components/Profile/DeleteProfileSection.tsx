import { useState } from 'react';
import { useAuth } from '@/hooks/useAuth';
import { deleteUser } from '@/lib/api';
import { useMutation } from '@tanstack/react-query';
import { Button } from '@/components/shared/Button';
import { ErrorMessage } from '@/components/shared/ErrorMessage';
import { IconTrash } from '@icons/CustomIcons';
import {
  Modal,
  ModalCard,
  ModalHeader,
  ModalTitle,
  ModalCloseButton,
  ModalBody,
  ModalFooter,
} from '@/components/Modal';
import styles from './DeleteProfileSection.module.css';

interface DeleteProfileSectionProps {
  onSuccess: () => void;
}

export const DeleteProfileSection = ({ onSuccess }: DeleteProfileSectionProps) => {
  const { login } = useAuth();
  const [isConfirmOpen, setIsConfirmOpen] = useState(false);

  const {
    mutate: performDelete,
    isPending: isDeleting,
    error: deleteError,
  } = useMutation({
    mutationFn: deleteUser,
    onSuccess: () => {
      login(null);
      setIsConfirmOpen(false);
      onSuccess();
    },
  });

  return (
    <div className={styles.sectionWrapper}>
      <div className={styles.dangerBox}>
        <div className={styles.dangerRow}>
          <div className={styles.iconWrapper}>
            <IconTrash />
          </div>
          <div className={styles.textGroup}>
            <h3 className={styles.title}>Danger Zone</h3>
            <p className={styles.description}>
              Permanently delete and anonymize your account.
            </p>
          </div>
          <div className={styles.btnWrapper}>
            <Button
              type="button"
              variant="destructive"
              onClick={() => setIsConfirmOpen(true)}
              disabled={isDeleting}
              data-testid="delete-profile-button"
              className={styles.deleteButton}
            >
              Delete Profile
            </Button>
          </div>
        </div>
      </div>

      {deleteError && (
        <div className={styles.errorContainer}>
          <ErrorMessage
            message={
              deleteError instanceof Error ? deleteError.message : 'Failed to delete profile'
            }
          />
        </div>
      )}

      {/* Custom Confirmation Modal */}
      <Modal isOpen={isConfirmOpen} onClose={() => setIsConfirmOpen(false)}>
        <ModalCard
          size="sm"
          className={styles.modalCard}
        >
          <ModalHeader className={styles.modalHeader}>
            <div className={styles.modalHeaderRow}>
              <IconTrash />
              <ModalTitle className={styles.modalTitle}>
                Confirm Deletion
              </ModalTitle>
            </div>
            <ModalCloseButton />
          </ModalHeader>

          <ModalBody className={styles.modalBody}>
            <p className={styles.modalText}>
              Are you sure you want to permanently delete and anonymize your profile?
            </p>
            <div className={styles.warningBox}>
              <span className={styles.warningTitle}>Warning</span>
              <span>
                This action cannot be undone. You will be logged out immediately and your username
                will be randomized.
              </span>
            </div>
          </ModalBody>

          <ModalFooter className={styles.modalFooter}>
            <Button
              type="button"
              variant="secondary"
              onClick={() => setIsConfirmOpen(false)}
              className={styles.footerButton}
              data-testid="delete-profile-confirm-cancel"
            >
              Cancel
            </Button>
            <Button
              type="button"
              variant="destructive"
              onClick={() => performDelete()}
              disabled={isDeleting}
              className={styles.confirmDeleteButton}
              data-testid="delete-profile-confirm-submit"
            >
              {isDeleting ? 'Deleting...' : 'Yes, Delete'}
            </Button>
          </ModalFooter>
        </ModalCard>
      </Modal>
    </div>
  );
};
