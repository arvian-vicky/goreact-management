import { colors, Paper, Typography } from '@mui/material';

import SidebarLayout from '@/components/layouts/SidebarLayout';

const Settings = () => {
  return (
    <SidebarLayout
      pageTitle="Pengaturan"
      breadcrumbs={[
        {
          label: 'Pengaturan',
        },
      ]}
    >
      <Paper
        sx={{
          padding: 2,
          background: colors.lightBlue[100],
        }}
      >
        <Typography>Menampilkan settings di sini</Typography>
      </Paper>
    </SidebarLayout>
  );
};

export default Settings;
