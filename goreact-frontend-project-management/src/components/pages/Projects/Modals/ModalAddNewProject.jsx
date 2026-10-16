import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';

import { CircularProgress, IconButton, Stack } from '@mui/material';
import { Close } from '@mui/icons-material';
import TextField from '@/components/ui/Forms/TextField';
import { useForm } from 'react-hook-form';
import { useState } from 'react';
import DatePicker from '@/components/ui/Forms/DatePicker';
import Modal from '@/components/ui/Modal';
import boards from '@/services/api/boards';
import lists from '@/services/api/lists';
import common from '@/utils/common';
import { useNavigate } from 'react-router';

const ModalAddNewProject = ({ open, handleClose }) => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const { control, handleSubmit, reset } = useForm();

  const onSubmit = async (data) => {
    setLoading(true);
    // create boards
    const response = await boards.create(data);

    // create lists
    const listTitles = ['Baru', 'Sedang dikerjakan', 'Selesai'];
    for (let i = 0; i < listTitles.length; i++) {
      await lists.create({
        board_public_id: response.data.data.public_id,
        title: listTitles[i],
      });
      await common.sleep(5000);
    }
    handleReset();
    navigate(`/projects/${response.data.data.public_id}`);
  };

  const handleReset = () => {
    reset();
    setLoading(false);
    handleClose();
  };

  return (
    <Modal
      open={open}
      handleClose={handleClose}
      aria-labelledby="modal-modal-title"
      aria-describedby="modal-modal-description"
      title={'Buat Proyek Baru'}
      sx={{
        minWidth: 400,
        maxWidth: 500,
      }}
    >
      <Box component={'form'} onSubmit={handleSubmit(onSubmit)}>
        {loading ? (
          <Box
            sx={{
              display: 'flex',
              justifyContent: 'center',
              alignItems: 'center',
              minHeight: 200,
            }}
          >
            <CircularProgress />
          </Box>
        ) : (
          <Stack
            sx={{
              p: 2,
            }}
          >
            <TextField label="Judul proyek" name="title" control={control} />
            <TextField
              label="Deskripsi"
              name="description"
              control={control}
              multiline
              rows={4}
              helperText={'Opsional'}
            />
            <DatePicker
              control={control}
              label="Tanggal Jatuh Tempo"
              name="due_date"
              sx={{
                width: '100%',
              }}
            />
          </Stack>
        )}

        <Stack
          direction="row"
          spacing={1}
          justifyContent="flex-end"
          sx={{ p: 2 }}
        >
          <Button
            type="submit"
            variant="contained"
            color="primary"
            disabled={loading}
          >
            Buat Proyek
          </Button>
          <Button
            type="button"
            variant="outlined"
            onClick={handleClose}
            disabled={loading}
          >
            Batal
          </Button>
        </Stack>
      </Box>
    </Modal>
  );
};

export default ModalAddNewProject;
