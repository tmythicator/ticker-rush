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

interface DeleteProfileSectionProps {
  onSuccess: () => void;
}

export const DeleteProfileSection = ({ onSuccess }: DeleteProfileSectionProps) => {
  const { login } = useAuth();
  const [isConfirmOpen, setIsConfirmOpen] = useState(false);

  const { mutate: performDelete, isPending: isDeleting, error: deleteError } = useMutation({
    mutationFn: deleteUser,
    onSuccess: () => {
      login(null);
      setIsConfirmOpen(false);
      onSuccess();
    },
  });

  return (
    <div className="space-y-4 pt-6 border-t border-border">
      <div className="rounded-xl border border-destructive/20 bg-destructive/5 p-4 sm:p-5 transition-all duration-200 hover:border-destructive/30">
        <div className="flex flex-col sm:flex-row sm:items-center gap-4">
          <div className="rounded-lg bg-destructive/10 p-2 text-destructive self-start sm:self-center">
            <IconTrash className="h-5 w-5" />
          </div>
          <div className="flex-1 space-y-0.5">
            <h3 className="text-sm font-semibold text-destructive">
              Danger Zone
            </h3>
            <p className="text-xs text-muted-foreground">
              Permanently delete and anonymize your account.
            </p>
          </div>
          <div className="self-end sm:self-center">
            <Button
              type="button"
              variant="destructive"
              onClick={() => setIsConfirmOpen(true)}
              disabled={isDeleting}
              data-testid="delete-profile-button"
              className="whitespace-nowrap transition-transform active:scale-95"
            >
              Delete Profile
            </Button>
          </div>
        </div>
      </div>

      {deleteError && (
        <div className="mt-2">
          <ErrorMessage
            message={deleteError instanceof Error ? deleteError.message : 'Failed to delete profile'}
          />
        </div>
      )}

      {/* Custom Confirmation Modal */}
      <Modal isOpen={isConfirmOpen} onClose={() => setIsConfirmOpen(false)}>
        <ModalCard size="sm" className="border border-destructive/30 shadow-2xl shadow-destructive/10">
          <ModalHeader className="mb-4">
            <div className="flex items-center gap-2 text-destructive">
              <IconTrash className="h-5 w-5" />
              <ModalTitle className="text-destructive font-bold text-lg">Confirm Deletion</ModalTitle>
            </div>
            <ModalCloseButton />
          </ModalHeader>

          <ModalBody className="space-y-4">
            <p className="text-sm text-foreground font-medium leading-relaxed">
              Are you sure you want to permanently delete and anonymize your profile?
            </p>
            <div className="rounded-lg bg-yellow-500/10 border border-yellow-500/20 p-3 text-xs text-yellow-500 space-y-1">
              <span className="font-semibold block">Warning</span>
              <span>This action cannot be undone. You will be logged out immediately and your username will be randomized.</span>
            </div>
          </ModalBody>

          <ModalFooter className="mt-6 gap-3">
            <Button
              type="button"
              variant="secondary"
              onClick={() => setIsConfirmOpen(false)}
              className="flex-1 justify-center"
              data-testid="delete-profile-confirm-cancel"
            >
              Cancel
            </Button>
            <Button
              type="button"
              variant="destructive"
              onClick={() => performDelete()}
              disabled={isDeleting}
              className="flex-1 justify-center shadow-lg shadow-destructive/20"
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
