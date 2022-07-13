package com.gusrylmubarok.health.backend.services;

import com.gusrylmubarok.health.backend.domain.Doctor;
import com.gusrylmubarok.health.backend.domain.User;

import java.util.List;

public interface DoctorService {
    void save(Doctor doctor);
    List<Doctor> findBySpeciality(String specialityCode);
    List<Doctor> findByLocation(String location);
    List<Doctor> findByHospital(String hospitalName);
    List<Doctor> findAll();
    Doctor findByUserEmailAddress(String email);
    int findCount();
    Doctor findByUserId(int userId);
    void addDoctor(User user);
}
