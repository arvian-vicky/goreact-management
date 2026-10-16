import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';

import { IconButton, Stack } from '@mui/material';
import Modal from '@/components/ui/Modal';
import { Close } from '@mui/icons-material';
import TextField from '@/components/ui/Forms/TextField';
import { useForm, useWatch } from 'react-hook-form';
import { useEffect, useState } from 'react';
import { useDebounce } from 'use-debounce';
import users from '@/services/api/users';

import * as Yup from 'yup';
import { yupResolver } from '@hookform/resolvers/yup';
import boards from '@/services/api/boards';
import { useParams } from 'react-router';

const searchSchema = Yup.object({
  search: Yup.string().required().email(),
});

const ModalAddNewMember = ({ open, handleClose }) => {
  const { id } = useParams();
  const {
    control,
    formState: { isValid },
    reset,
  } = useForm({
    resolver: yupResolver(searchSchema),
    defaultValues: {
      search: '',
    },
  });
  const [members, setMembers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [hasError, setHasError] = useState();

  const watchSearch = useWatch({
    control,
    name: 'search',
  });
  const [debounceSearch] = useDebounce(watchSearch, 1000);

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        setLoading(true);
        const response = await users.findAll({
          filter: debounceSearch,
        });
        setMembers(response.data.data);
      } catch (error) {
        setHasError(error?.response?.data);
      } finally {
        setLoading(false);
      }
    };

    if (debounceSearch) {
      fetchUsers();
    }
  }, [debounceSearch]);

  const handleAddNewMember = async () => {
    setLoading(true);
    await boards.addMembers(id, members[0]?.public_id);
    setLoading(false);
    reset();
    handleClose();
  };

  return (
    <Modal
      open={open}
      handleClose={handleClose}
      aria-labelledby="modal-modal-title"
      aria-describedby="modal-modal-description"
      title={`Menambahkan anggota pada proyek ini.`}
      sx={{
        minWidth: 800,
        maxWidth: 1000,
      }}
    >
      <Stack
        sx={{
          p: 2,
        }}
      >
        <Box width={'100%'}>
          <TextField
            control={control}
            name={'search'}
            label={'Cari email pengguna'}
            placeholder="Masukkan email pengguna"
            error={hasError?.message ?? ''}
            helperText={hasError?.message ?? ''}
          />
        </Box>
        <Stack justifyContent={'flex-end'} alignItems={'flex-end'}>
          <Button
            type="button"
            variant="contained"
            disabled={!isValid || hasError}
            loading={loading}
            onClick={handleAddNewMember}
          >
            Tambahkan
          </Button>
        </Stack>
      </Stack>
    </Modal>
  );
};

export default ModalAddNewMember;
