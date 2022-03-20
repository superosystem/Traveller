package com.gusrylmubarok.springdi;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.AnnotationConfigApplicationContext;


public class SpringDiApplication {

	public static void main(String[] args) {
		AnnotationConfigApplicationContext context = new AnnotationConfigApplicationContext(AppConfig.class);
		FirstBean bean = context.getBean(FirstBean.class);
		bean.populateBeans();
		context.close();
	}

}
