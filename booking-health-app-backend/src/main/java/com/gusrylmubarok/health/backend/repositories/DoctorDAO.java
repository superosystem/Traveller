package com.gusrylmubarok.health.backend.repositories;

import com.gusrylmubarok.health.backend.domain.Doctor;

import java.util.List;

public interface DoctorDAO {
    List<Doctor> findAll();
    List<Doctor> findBySpecialityCode(String code);
    int findAllCount();
    List<Doctor> findByLocation(String location);
    List<Doctor> findByHospital(String hospitalName);
    Doctor findByUserId(int userId);
    Doctor save(Doctor doctor);
}
