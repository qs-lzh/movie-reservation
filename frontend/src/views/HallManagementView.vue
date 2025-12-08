<template>
  <div class="hall-management-container">
    <h2>Hall Management</h2>
    
    <!-- Add Hall Form -->
    <el-card class="add-hall-card">
      <template #header>
        <span>Add New Hall</span>
      </template>
      <el-form :model="hallForm" :rules="formRules" ref="hallFormRef" label-width="120px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Name" prop="name">
              <el-input v-model="hallForm.name" placeholder="Enter hall name"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Seat Count" prop="seatCount">
              <el-input-number 
                v-model="hallForm.seatCount" 
                :min="1" 
                placeholder="Enter seat count"
                style="width: 100%"></el-input-number>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Rows" prop="rows">
              <el-input-number 
                v-model="hallForm.rows" 
                :min="1" 
                placeholder="Enter number of rows"
                style="width: 100%"></el-input-number>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Columns" prop="cols">
              <el-input-number 
                v-model="hallForm.cols" 
                :min="1" 
                placeholder="Enter number of columns"
                style="width: 100%"></el-input-number>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item>
          <el-button type="primary" @click="handleAddHall" :loading="isSubmitting">Add Hall</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- Halls Table -->
    <el-card class="halls-table-card">
      <template #header>
        <span>All Halls</span>
      </template>
      <el-table :data="halls" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="100"></el-table-column>
        <el-table-column prop="name" label="Name"></el-table-column>
        <el-table-column prop="seat_count" label="Seat Count" width="120"></el-table-column>
        <el-table-column prop="rows" label="Rows" width="100"></el-table-column>
        <el-table-column prop="cols" label="Columns" width="120"></el-table-column>
        <el-table-column label="Actions" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">Edit</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row.id)">Delete</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Edit Hall Dialog -->
    <el-dialog v-model="editDialogVisible" title="Edit Hall" width="500px">
      <el-form :model="editForm" :rules="formRules" ref="editFormRef" label-width="120px">
        <el-form-item label="Name" prop="name">
          <el-input v-model="editForm.name" placeholder="Enter hall name"></el-input>
        </el-form-item>
        <el-form-item label="Seat Count" prop="seatCount">
          <el-input-number 
            v-model="editForm.seatCount" 
            :min="1" 
            placeholder="Enter seat count"
            style="width: 100%"></el-input-number>
        </el-form-item>
        <el-form-item label="Rows" prop="rows">
          <el-input-number 
            v-model="editForm.rows" 
            :min="1" 
            placeholder="Enter number of rows"
            style="width: 100%"></el-input-number>
        </el-form-item>
        <el-form-item label="Columns" prop="cols">
          <el-input-number 
            v-model="editForm.cols" 
            :min="1" 
            placeholder="Enter number of columns"
            style="width: 100%"></el-input-number>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="handleUpdateHall" :loading="isUpdating">Update</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { hallAPI } from '@/api/hall'

const halls = ref([])
const loading = ref(false)
const isSubmitting = ref(false)
const isUpdating = ref(false)

// Add hall form
const hallForm = ref({
  name: '',
  seatCount: 100,
  rows: 10,
  cols: 10
})

// Edit hall form
const editForm = ref({
  id: null,
  name: '',
  seatCount: 100,
  rows: 10,
  cols: 10
})

const editDialogVisible = ref(false)
const hallFormRef = ref()
const editFormRef = ref()

const formRules = {
  name: [
    { required: true, message: 'Please enter hall name', trigger: 'blur' },
    { min: 1, max: 100, message: 'Hall name must be between 1 and 100 characters', trigger: 'blur' }
  ],
  seatCount: [
    { required: true, message: 'Please enter seat count', trigger: 'blur' },
    { type: 'number', min: 1, message: 'Seat count must be at least 1', trigger: 'blur' }
  ],
  rows: [
    { required: true, message: 'Please enter number of rows', trigger: 'blur' },
    { type: 'number', min: 1, message: 'Rows must be at least 1', trigger: 'blur' }
  ],
  cols: [
    { required: true, message: 'Please enter number of columns', trigger: 'blur' },
    { type: 'number', min: 1, message: 'Columns must be at least 1', trigger: 'blur' }
  ]
}

const fetchHalls = async () => {
  loading.value = true
  try {
    const response = await hallAPI.getHalls()
    halls.value = response.data.data
  } catch (error) {
    console.error('Failed to fetch halls:', error)
    ElMessage.error('Failed to fetch halls')
  } finally {
    loading.value = false
  }
}

const handleAddHall = async () => {
  if (!hallFormRef.value) return
  
  await hallFormRef.value.validate(async (valid) => {
    if (valid) {
      isSubmitting.value = true
      try {
        await hallAPI.createHall(
          hallForm.value.name,
          hallForm.value.seatCount,
          hallForm.value.rows,
          hallForm.value.cols
        )
        ElMessage.success('Hall created successfully')
        
        // Reset form
        hallForm.value = {
          name: '',
          seatCount: 100,
          rows: 10,
          cols: 10
        }
        
        // Refresh halls list
        await fetchHalls()
      } catch (error) {
        console.error('Failed to create hall:', error)
        ElMessage.error('Failed to create hall')
      } finally {
        isSubmitting.value = false
      }
    }
  })
}

const handleEdit = (row) => {
  editForm.value = {
    id: row.id,
    name: row.name,
    seatCount: row.seat_count,
    rows: row.rows,
    cols: row.cols
  }
  editDialogVisible.value = true
}

const handleUpdateHall = async () => {
  if (!editFormRef.value) return
  
  await editFormRef.value.validate(async (valid) => {
    if (valid) {
      isUpdating.value = true
      try {
        await hallAPI.updateHall(
          editForm.value.id,
          editForm.value.name,
          editForm.value.seatCount,
          editForm.value.rows,
          editForm.value.cols
        )
        ElMessage.success('Hall updated successfully')
        editDialogVisible.value = false
        
        // Refresh halls list
        await fetchHalls()
      } catch (error) {
        console.error('Failed to update hall:', error)
        ElMessage.error('Failed to update hall')
      } finally {
        isUpdating.value = false
      }
    }
  })
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm(
      'Are you sure you want to delete this hall?',
      'Warning',
      {
        confirmButtonText: 'OK',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }
    )
    
    await hallAPI.deleteHall(id)
    ElMessage.success('Hall deleted successfully')
    
    // Refresh halls list
    await fetchHalls()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete hall:', error)
      ElMessage.error('Failed to delete hall')
    }
  }
}

onMounted(() => {
  fetchHalls()
})
</script>

<style scoped>
.hall-management-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.add-hall-card {
  margin-bottom: 20px;
}

.halls-table-card {
  margin-bottom: 20px;
}
</style>