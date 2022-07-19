package com.gusrylmubarok.spring.perpustakaan;

import com.gusrylmubarok.spring.perpustakaan.entity.User;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.security.core.context.SecurityContextHolder;

@SpringBootApplication
public class PerpustakaanSpringRestApplication {
	public static void main(String[] args) {
		SpringApplication.run(PerpustakaanSpringRestApplication.class, args);
	}

	public static User getCurrentUser() {
		try {
			Object principal = SecurityContextHolder.getContext().getAuthentication()
					.getPrincipal();
			if (principal != null && principal.getClass().equals(User.class)) {
				return (User) principal;
			}
		}catch (Exception ignore) {
		}
		return null;
	}
}
