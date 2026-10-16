import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';

import { IconButton, Stack } from '@mui/material';
import Modal from '@/components/ui/Modal';
import { Close } from '@mui/icons-material';

const ModalCloseProject = ({ open, handleClose }) => {
  return (
    <Modal
      open={open}
      onClose={handleClose}
      aria-labelledby="modal-modal-title"
      aria-describedby="modal-modal-description"
      title={'Apakah Anda yakin ingin menutup proyek ini?'}
      sx={{
        maxWidth: 300,
      }}
    >
      <Box
        sx={{
          p: 2,
        }}
      >
        <Typography id="modal-modal-description" sx={{ mt: 2 }}>
          Setelah proyek ditutup, Anda tidak akan dapat menambahkan atau
          mengedit tugas lagi.
        </Typography>
      </Box>
      <Stack direction="row" spacing={1} p={2} justifyContent="flex-end" mt={3}>
        <Button
          variant="contained"
          color="error"
          onClick={handleClose}
          sx={{ mt: 2 }}
        >
          Tutup Proyek
        </Button>
        <Button
          variant="outlined"
          onClick={handleClose}
          sx={{ mt: 2, ml: 1 }}
          color="primary"
        >
          Batal
        </Button>
      </Stack>
    </Modal>
  );
};

export default ModalCloseProject;
