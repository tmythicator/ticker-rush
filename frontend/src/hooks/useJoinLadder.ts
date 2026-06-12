import { useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { joinLadder } from '@/lib/api';
import { QUERY_KEY_USER } from '@/lib/queryKeys';

export const useJoinLadder = () => {
  const queryClient = useQueryClient();
  const [isConfirming, setIsConfirming] = useState(false);

  const mutation = useMutation({
    mutationFn: joinLadder,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: QUERY_KEY_USER });
    },
  });

  const handleJoin = () => {
    mutation.mutate();
  };

  return {
    isConfirming,
    setIsConfirming,
    isPending: mutation.isPending,
    handleJoin,
  };
};
